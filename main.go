// iglu is a small program for generating Monero cold wallets.
//
// The basis of this implementation is mostly explained at
// https://xmr.llcoins.net/addresstests.html.
package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

	"filippo.io/edwards25519"
	"golang.org/x/crypto/sha3"
)

var l *big.Int

func init() {
	var ok bool
	l, ok = big.NewInt(0).SetString("27742317777372353535851937790883648493", 10)
	if !ok {
		panic("could not create big integer")
	}
	l.SetBit(l, 252, 1)
}

func main() {
	secretSpendKey, err := generateSecretSpendKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while generating secret spend key: %v\n", err)
		os.Exit(1)
	}
	secretViewKey, err := deriveSecretViewKey(secretSpendKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while hashing secret spend key: %v\n", err)
		os.Exit(1)
	}
	publicSpendKey, err := derivePublicKey(secretSpendKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while deriving public spend key: %v\n", err)
		os.Exit(1)
	}
	publicViewKey, err := derivePublicKey(secretViewKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while deriving public view key: %v\n", err)
		os.Exit(1)
	}
	address, err := deriveAddress(publicSpendKey, publicViewKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while deriving address: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("secret spend key: %s\n", toString(secretSpendKey))
	fmt.Printf("secret view key : %s\n", toString(secretViewKey))
	fmt.Printf("primary address : %s\n", address)
}

func generateSecretSpendKey() (*big.Int, error) {
	// Theoretically, 0 would have to be prevented.
	key, err := rand.Int(rand.Reader, l)
	if err != nil {
		return key, err
	}
	return revertBytes(key), nil
}

func deriveSecretViewKey(secretSpendKey *big.Int) (*big.Int, error) {
	bytes := make([]byte, 32, 32)
	secretSpendKey.FillBytes(bytes)
	hasher := sha3.NewLegacyKeccak256()
	if _, err := hasher.Write(bytes); err != nil {
		return big.NewInt(0), err
	}
	secretViewKey := big.NewInt(0).SetBytes(hasher.Sum(nil))
	secretViewKey = revertBytes(secretViewKey)
	secretViewKey.Mod(secretViewKey, l)
	return revertBytes(secretViewKey), nil
}

func derivePublicKey(secretKey *big.Int) (*big.Int, error) {
	bytes := make([]byte, 32, 32)
	secretKey.FillBytes(bytes)
	s, err := edwards25519.NewScalar().SetCanonicalBytes(bytes)
	if err != nil {
		return big.NewInt(0), err
	}
	p := edwards25519.NewIdentityPoint().ScalarBaseMult(s)
	return big.NewInt(0).SetBytes(p.Bytes()), nil
}

func deriveAddress(publicSpendKey, publicViewKey *big.Int) (string, error) {
	publicSpendKeyBytes := make([]byte, 32, 32)
	publicViewKeyBytes := make([]byte, 32, 32)
	publicSpendKey.FillBytes(publicSpendKeyBytes)
	publicViewKey.FillBytes(publicViewKeyBytes)
	bytes := make([]byte, 0, 69)
	bytes = append(bytes, 0x12) // Monero network byte.
	bytes = append(bytes, publicSpendKeyBytes...)
	bytes = append(bytes, publicViewKeyBytes...)
	hasher := sha3.NewLegacyKeccak256()
	if _, err := hasher.Write(bytes); err != nil {
		return "", err
	}
	hashBytes := hasher.Sum(nil)
	bytes = append(bytes, hashBytes[:4]...)
	return moneroAddressEncode(bytes), nil
}

func moneroAddressEncode(bytes []byte) (out string) {
	for i := 0; i < 8; i++ {
		block := base58Encode(bytes[8*i : 8*(i+1)])
		out += strings.Repeat("1", 11-len(block)) + block
	}
	block := base58Encode(bytes[64:69])
	return out + strings.Repeat("1", 7-len(block)) + block
}

// base58Encode encodes bytes as base58 using simplified code from
// https://github.com/akamensky/base58.
func base58Encode(in []byte) string {
	alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	out := make([]byte, 0)
	zero := big.NewInt(0)
	bigRadix := big.NewInt(58)
	num := big.NewInt(0).SetBytes(in)
	for num.Cmp(zero) > 0 {
		mod := big.NewInt(0)
		num.DivMod(num, bigRadix, mod)
		out = append(out, alphabet[mod.Int64()])
	}
	for i := 0; i < len(out)/2; i++ {
		out[i], out[len(out)-1-i] = out[len(out)-1-i], out[i]
	}
	return string(out)
}

// revertBytes reverts the bytes of a 32byte *big.Int and returns it.
// This is used to change the endianness of a *big.Int.
func revertBytes(x *big.Int) *big.Int {
	bytes := make([]byte, 32, 32)
	x.FillBytes(bytes)
	for i := 0; i < 16; i++ {
		bytes[i], bytes[31-i] = bytes[31-i], bytes[i]
	}
	return x.SetBytes(bytes)
}

// toString converts a key into its 32bit hexadecimal representation.
func toString(x *big.Int) string {
	bytes := make([]byte, 32, 32)
	x.FillBytes(bytes)
	return fmt.Sprintf("%x", bytes)
}
