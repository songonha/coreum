package cored

// well-known keys to create predictable wallets so manual operation is easier
var (
	// AlicePrivKey is the private key of alice having address: cosmos1x645ym2yz4gckqjtpwr8yddqzkkzdpktxr3clr
	AlicePrivKey = Secp256k1PrivateKey{0x9b, 0xc9, 0xd0, 0x15, 0x11, 0x2, 0x94, 0x9, 0x92, 0xfd, 0x2b, 0xad, 0xbe, 0x36, 0x63, 0x1f, 0xed, 0x30, 0x10, 0xd, 0x6e, 0x24, 0xb1, 0xc2, 0x58, 0xb4, 0xfd, 0xe4, 0xaf, 0xdd, 0xf2, 0x40}
	// BobPrivKey is the private key of bob having address: cosmos1cjs7qela0trw2qyyfxw5e5e7cvwzprkju5d2su
	BobPrivKey = Secp256k1PrivateKey{0x87, 0x70, 0x0, 0x22, 0xa3, 0x24, 0x81, 0x59, 0x8d, 0xb8, 0x27, 0x57, 0xdb, 0x97, 0xe6, 0x9b, 0xed, 0x11, 0xb6, 0x17, 0x3, 0xcc, 0x44, 0xe0, 0x2a, 0xd3, 0x1e, 0x95, 0x36, 0xcf, 0x2d, 0x7f}
	// CharliePrivKey is the private key of charlie having address: cosmos1rd8wynz2987ey6pwmkuwfg9q8hf04xdyjqy2f4
	CharliePrivKey = Secp256k1PrivateKey{0x12, 0x9, 0x56, 0x3d, 0x40, 0x69, 0xf7, 0x57, 0xdd, 0x4c, 0x69, 0x17, 0x92, 0x7, 0xf0, 0xe6, 0x62, 0xa1, 0xcb, 0x8c, 0xfe, 0x8, 0x61, 0x68, 0x4c, 0x5e, 0xbc, 0x6b, 0x34, 0xa9, 0x5f, 0x7}
)