package zstress

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"runtime"

	"github.com/CoreumFoundation/coreum-tools/pkg/logger"
	"github.com/CoreumFoundation/coreum-tools/pkg/must"
	"github.com/CoreumFoundation/coreum-tools/pkg/parallel"
	"go.uber.org/zap"

	"github.com/CoreumFoundation/coreum/coreznet/infra/apps/cored"
)

// StressConfig contains config for benchmarking the blockchain
type StressConfig struct {
	// ChainID is the ID of the chain to connect to
	ChainID string

	// NodeAddress is the address of a cored node RPC endpoint, in the form of host:port, to connect to
	NodeAddress string

	// Accounts is the list of private keys used to send transactions during benchmark
	Accounts []cored.Secp256k1PrivateKey

	// NumOfAccounts is the number of accounts used to benchmark the node in parallel
	NumOfAccounts int

	// NumOfTransactions to send from each account
	NumOfTransactions int
}

type tx struct {
	AccountIndex int
	TxIndex      int
	From         cored.Wallet
	To           cored.Wallet
	TxBytes      []byte
}

// Stress runs a benchmark test
func Stress(ctx context.Context, config StressConfig) error {
	if config.NumOfAccounts <= 0 {
		return errors.New("number of accounts must be greater than 0")
	}
	if config.NumOfAccounts > len(config.Accounts) {
		return errors.New("number of accounts is greater than the number of provided private keys")
	}

	log := logger.Get(ctx)
	client := cored.NewClient(config.ChainID, config.NodeAddress)

	log.Info("Preparing signed transactions...")
	var signedTxs [][][]byte
	err := parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		queue := make(chan tx)
		results := make(chan tx)
		for i := 0; i < runtime.NumCPU(); i++ {
			spawn(fmt.Sprintf("signer-%d", i), parallel.Continue, func(ctx context.Context) error {
				for {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case tx, ok := <-queue:
						if !ok {
							return nil
						}
						tx.TxBytes = must.Bytes(client.TxBankSend(tx.From, tx.To, cored.Balance{Amount: big.NewInt(1), Denom: "core"}))
						select {
						case <-ctx.Done():
							return ctx.Err()
						case results <- tx:
						}
					}
				}
			})
		}
		spawn("enqueue", parallel.Continue, func(ctx context.Context) error {
			for i := 0; i < config.NumOfAccounts; i++ {
				fromPrivateKey := config.Accounts[i]
				toPrivateKeyIndex := i + 1
				if toPrivateKeyIndex >= config.NumOfAccounts {
					toPrivateKeyIndex = 0
				}
				toPrivateKey := config.Accounts[toPrivateKeyIndex]

				accNum, accSeq, err := client.GetNumberSequence(fromPrivateKey.Address())
				if err != nil {
					return fmt.Errorf("fetching account number and sequence failed: %w", err)
				}

				tx := tx{
					AccountIndex: i,
					From:         cored.Wallet{Name: "sender", Key: fromPrivateKey, AccountNumber: accNum, AccountSequence: accSeq},
					To:           cored.Wallet{Name: "receiver", Key: toPrivateKey},
				}

				for j := 0; j < config.NumOfTransactions; j++ {
					tx.TxIndex = j
					select {
					case <-ctx.Done():
						return ctx.Err()
					case queue <- tx:
					}
					tx.From.AccountSequence++
				}
			}
			return nil
		})
		spawn("integrate", parallel.Exit, func(ctx context.Context) error {
			signedTxs = make([][][]byte, config.NumOfAccounts)
			for i := 0; i < config.NumOfAccounts; i++ {
				signedTxs[i] = make([][]byte, config.NumOfTransactions)
			}
			for i := 0; i < config.NumOfAccounts; i++ {
				for j := 0; j < config.NumOfTransactions; j++ {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case result := <-results:
						signedTxs[result.AccountIndex][result.TxIndex] = result.TxBytes
					}
				}
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return err
	}
	log.Info("Transactions prepared")

	log.Info("Broadcasting transactions...")
	err = parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		for i, accountTxs := range signedTxs {
			accountTxs := accountTxs
			spawn(fmt.Sprintf("account-%d", i), parallel.Continue, func(ctx context.Context) error {
				log := logger.Get(ctx)
				j := 0
				for {
					if err := ctx.Err(); err != nil {
						return err
					}
					tx := accountTxs[j]
					txHash, err := client.Broadcast(tx)
					if err != nil {
						log.Error("Sending transaction failed", zap.Error(err))
						continue
					}
					log.Debug("Transaction broadcasted", zap.String("txHash", txHash))

					j++
					if j >= config.NumOfTransactions {
						return nil
					}
				}
			})
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Info("Benchmark finished")
	return nil
}