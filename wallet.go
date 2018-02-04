package main

/**
Reference below link for more detail
https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
*/

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const publicKeyPrefix = byte(0x04)
const networkVersion = byte(0x00)
const addressChecksumLen = 4

// Wallet contains two important keys which make its identification
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func newWallet() *Wallet {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKeyFirst33Bytes := append([]byte{publicKeyPrefix}, privateKey.PublicKey.X.Bytes()...)
	publicKey := append(publicKeyFirst33Bytes, privateKey.PublicKey.Y.Bytes()...)

	return &Wallet{*privateKey, publicKey}
}

func (w Wallet) getAddress() string {
	publicKeyHash := hashPublicKey(w.PublicKey)

	versionPayload := append([]byte{networkVersion}, publicKeyHash...)
	checksum := publickeyChecksum(versionPayload)

	fullPayload := append(versionPayload, checksum...)
	address := base58.Encode(fullPayload)

	return address
}

func validateAddress(address string) bool {
	payload := base58.Decode(address)
	actualChecksum := payload[len(payload)-addressChecksumLen:]
	version := payload[0]
	publicKeyHash := payload[1 : len(payload)-addressChecksumLen]
	targetChecksum := publickeyChecksum(append([]byte{version}, publicKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func hashPublicKey(publicKey []byte) []byte {
	publicSHA256 := sha256.Sum256(publicKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEDEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEDEMD160
}

func publickeyChecksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
