// iglu is a small program for generating Monero cold wallets.
//
// The basis of this implementation is mostly explained at
// https://xmr.llcoins.net/addresstests.html.
package main

import (
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"filippo.io/edwards25519"
	"golang.org/x/crypto/sha3"
)

var l *big.Int
var nSubadresses = flag.Int("s", 0,
	"Number of subaddresses to generate. The account index is always 0.")

func init() {
	var ok bool
	l, ok = big.NewInt(0).SetString("27742317777372353535851937790883648493", 10)
	if !ok {
		panic("could not create big integer")
	}
	l.SetBit(l, 252, 1)
}

func main() {
	flag.Parse()
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
	address, err := deriveAddress(publicSpendKey, publicViewKey, 18)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while deriving address: %v\n", err)
		os.Exit(1)
	}
	subaddresses := make([]string, 0)
	// Subaddress 0 is also skipped in the official wallet:
	for i := 1; nSubadresses != nil && i <= *nSubadresses; i++ {
		subaddress, err := deriveSubaddress(secretViewKey, publicSpendKey, uint32(0), uint32(i))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while deriving subaddress: %v\n", err)
			os.Exit(1)
		}
		subaddresses = append(subaddresses, subaddress)
	}

	fmt.Printf("secret spend key: %s\n", toString(secretSpendKey))
	fmt.Printf("secret view key : %s\n", toString(secretViewKey))
	fmt.Printf("primary address : %s\n", address)
	for i, subaddress := range subaddresses {
		fmt.Printf("subaddress #%04d: %s\n", i+1, subaddress)
	}
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

func deriveAddress(publicSpendKey, publicViewKey *big.Int, networkByte byte) (string, error) {
	publicSpendKeyBytes := make([]byte, 32, 32)
	publicViewKeyBytes := make([]byte, 32, 32)
	publicSpendKey.FillBytes(publicSpendKeyBytes)
	publicViewKey.FillBytes(publicViewKeyBytes)
	bytes := make([]byte, 0, 69)
	bytes = append(bytes, networkByte)
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

// deriveSubaddress derives the subaddress for account iAcc and index
// iSub for the standard address' secretViewKey and publicSpendKey.
// See https://monerodocs.org/public-address/subaddress/.
func deriveSubaddress(secretViewKey, publicSpendKey *big.Int, iAcc, iSub uint32) (string, error) {
	m, err := deriveM(secretViewKey, iAcc, iSub)
	if err != nil {
		return "", fmt.Errorf("could not derive m: %w", err)
	}
	publicSpendKeyBytes := make([]byte, 32, 32)
	publicSpendKey.FillBytes(publicSpendKeyBytes)
	B, err := edwards25519.NewIdentityPoint().SetBytes(publicSpendKeyBytes)
	if err != nil {
		return "", err
	}
	D := B.Add(B, edwards25519.NewIdentityPoint().ScalarBaseMult(m))
	subPublicSpendKey := big.NewInt(0).SetBytes(D.Bytes())

	secretViewKeyBytes := make([]byte, 32, 32)
	secretViewKey.FillBytes(secretViewKeyBytes)
	a, err := edwards25519.NewScalar().SetCanonicalBytes(secretViewKeyBytes)
	if err != nil {
		return "", err
	}
	C := edwards25519.NewIdentityPoint().ScalarMult(a, D)
	subPublicViewKey := big.NewInt(0).SetBytes(C.Bytes())

	return deriveAddress(subPublicSpendKey, subPublicViewKey, 42)
}

func deriveM(secretViewKey *big.Int, iAcc, iSub uint32) (*edwards25519.Scalar, error) {
	hasher := sha3.NewLegacyKeccak256()
	if _, err := hasher.Write([]byte("SubAddr\u0000")); err != nil {
		return edwards25519.NewScalar(), err
	}
	secretViewKeyBytes := make([]byte, 32, 32)
	secretViewKey.FillBytes(secretViewKeyBytes)
	if _, err := hasher.Write(secretViewKeyBytes); err != nil {
		return edwards25519.NewScalar(), err
	}
	buf := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buf, iAcc)
	if _, err := hasher.Write(buf); err != nil {
		return edwards25519.NewScalar(), err
	}
	buf = make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buf, iSub)
	if _, err := hasher.Write(buf); err != nil {
		return edwards25519.NewScalar(), err
	}
	mInt := big.NewInt(0).SetBytes(hasher.Sum(nil))
	revertBytes(mInt)
	mInt.Mod(mInt, l)
	revertBytes(mInt)
	mBytes := make([]byte, 32, 32)
	mInt.FillBytes(mBytes)
	return edwards25519.NewScalar().SetCanonicalBytes(mBytes)
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

// toString converts a key into its 32byte hexadecimal representation.
func toString(x *big.Int) string {
	bytes := make([]byte, 32, 32)
	x.FillBytes(bytes)
	return fmt.Sprintf("%x", bytes)
}
