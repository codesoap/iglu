package main

import (
	"math/big"
	"testing"
)

// testCases are some keys and addresses I generated with the official
// monero-wallet-cli.
var testCases = []struct {
	secretSpendKey string
	secretViewKey  string
	address        string
	subaddress1    string
	subaddress2    string
	subaddress3    string
}{
	{
		"a51ac194dae681f2d536546e6177befaf62c37beccead34e49f5dd63a057050a",
		"97a5c64e3fb463819da00459da06b10b4b9284bbede3f7af836b3860e4d43803",
		"469tuBh3CzsJXkRNvYEiZrWkHZ2xq8BeeaaBH7qjT9jV5KCBYMjfXuR1bSiXeKJTgyRPWywTkYr5PdvY4az9e8z1S5cuuvY",
		"84VDMTwfjPdXoSUqFWZoWciohBa4nXvw2L6eUvpN12YjUVuyAjaGiFNLNBoSCZfh1hMgZD4o4AX74V9u1gvcaDH4T61jk1w",
		"84LUo37osmTYsjma8VEVT6QfkieLyPr1mRtfhcxqenhp8ZEnaqkx4apgtsZmGUmVUzbehZYo7FNt6P2ysXAfDvKdBUQWjmd",
		"832B6JKFqiqH7MepXMTpNy6KVuZ5wkaBiPxkWGgsdbU7enya5B8QWTTEtTSU9UXWtLaEvy1dXZkxoKPYTB7FfS1NJNksYD8",
	},
	{
		"8dfb1ea20c0df41de5077db9bd3b35a461616175f4cbaa833a82865ae0079000",
		"bce9aa4d091c6addb224fcc57ed8ba96dc2acade62f25bfe0b64aa8b9dca710c",
		"48nXrUSR57zinEBVphSb7tRhj5qtakUQaaVRW6DJihqxQvn28rZRJvaNTzm9Ku2rpRKhqwLovm6aFDYYfBvBDmCKVCuXSpq",
		"82jV7SMB8EVEVJRhuv8jKLhhgaA3tnVPb6XD6mV2YegBgp68PZnFna1Sfj5qXMio5MjX5Gu9RYUaRWYGP2i52StwMeSypYg",
		"84Co5qp661WMzaGiUmFvqk3RRV2qWYSRufwtQAhD2BXQJTQB57btspVWUKfHoipgpZKJ9ktfbLvFPJKu93zqjnADNqJuxZy",
		"8AutHMLSyzGdTkVCQNLTUZQMeVKxgFyQ4XBpKq7y2qc7hruPTj3BPXgKNVUykAcX4FJte6KXJH5WSeJSxMPqqoRdM7uHasc",
	},
	{
		"05208f857fa500fe47a9377a13fa0af93374b1e245f12c46907e740b8b686d07",
		"d8a1c8d580f37009b65b78a55b87043fd96c4e4e908fa3471707a37cbf8c210a",
		"451ePSa7A2rBTB3fuR5bodXpty13nHanE3MoDBZEcdM7V6TgaSvQvd68ewNb9NcW4e6XTyiyi4mrq5S9gK6MZqfg9p6cP73",
		"8BGPdSGTazEeXfV1ksJt7vidMU936KarAVEPRmbx9XG8PJm6fQ6rEYDYWxNK6H5EBfKJQRJLw9s722s9Li1bNNXPEgSVbgk",
		"88NjEWAdGa1WmLpjEVNfpP47FcXkAAA8yJSDk7m57KkzB4dGXndZnkm3RbtApYCiUgHpEQSYRu7vSdHFhGgcxNPD2wRHqR5",
		"8BfU2j7SqyAaY48x3PRQEde6vG1ej3Dvrdt1sdM1ABMvPJXpYLUz1NHEY3jncv9M1P7iXgGx8fPgEZ6nrbpyKVLHDPtQRHr",
	},
}

func TestAll(t *testing.T) {
	for _, testCase := range testCases {
		secretSpendKey, ok := big.NewInt(0).SetString(testCase.secretSpendKey, 16)
		if !ok {
			t.Errorf("Could not parse secret spend key: %s.", testCase.secretSpendKey)
		}

		secretViewKey, err := deriveSecretViewKey(secretSpendKey)
		if err != nil {
			t.Errorf("Error while hashing secret spend key: %v", err)
		}
		if toString(secretViewKey) != testCase.secretViewKey {
			t.Errorf("Wrong secret view key: want %s, got: %s.", testCase.secretViewKey, toString(secretViewKey))
		}

		publicSpendKey, err := derivePublicKey(secretSpendKey)
		if err != nil {
			t.Errorf("Error while deriving public spend key: %v", err)
		}
		publicViewKey, err := derivePublicKey(secretViewKey)
		if err != nil {
			t.Errorf("Error while deriving public view key: %v", err)
		}

		address, err := deriveAddress(publicSpendKey, publicViewKey, 18)
		if err != nil {
			t.Errorf("Error while deriving address: %v", err)
		}
		if address != testCase.address {
			t.Errorf("Wrong address: want %s, got: %s.", testCase.address, address)
		}

		subaddress1, err := deriveSubaddress(secretViewKey, publicSpendKey, uint32(0), uint32(1))
		if err != nil {
			t.Errorf("Error while deriving subaddress1: %v", err)
		}
		if subaddress1 != testCase.subaddress1 {
			t.Errorf("Wrong subaddress1: want %s, got: %s.", testCase.subaddress1, subaddress1)
		}

		subaddress2, err := deriveSubaddress(secretViewKey, publicSpendKey, uint32(0), uint32(2))
		if err != nil {
			t.Errorf("Error while deriving subaddress2: %v", err)
		}
		if subaddress2 != testCase.subaddress2 {
			t.Errorf("Wrong subaddress2: want %s, got: %s.", testCase.subaddress2, subaddress2)
		}

		subaddress3, err := deriveSubaddress(secretViewKey, publicSpendKey, uint32(0), uint32(3))
		if err != nil {
			t.Errorf("Error while deriving subaddress3: %v", err)
		}
		if subaddress3 != testCase.subaddress3 {
			t.Errorf("Wrong subaddress3: want %s, got: %s.", testCase.subaddress3, subaddress3)
		}
	}
}
