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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

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
	Address    string
}

// StorableWallet is brief type of Wallet, which can be stored to paper
type StorableWallet struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

var wallet *Wallet

func setWallet(w *Wallet) {
	wallet = w
}

func getWallet() *Wallet {
	return wallet
}

func newWallet() *Wallet {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKeyFirst33Bytes := append([]byte{publicKeyPrefix}, privateKey.PublicKey.X.Bytes()...)
	publicKey := append(publicKeyFirst33Bytes, privateKey.PublicKey.Y.Bytes()...)
	publicKeyHash := hashPublicKey(publicKey)
	address := generateAddress(publicKeyHash)

	return &Wallet{*privateKey, publicKey, address}
}

func generateAddress(publicKeyHash []byte) string {

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

func (w *Wallet) toStorable() *StorableWallet {
	sWallet := new(StorableWallet)
	sWallet.PrivateKey = hex.EncodeToString(w.PrivateKey.D.Bytes())
	sWallet.PublicKey = hex.EncodeToString(w.PublicKey)
	sWallet.Address = w.Address
	return sWallet
}

func (sW *StorableWallet) toWallet() *Wallet {
	w := new(Wallet)
	curve := elliptic.P256()
	privateKeyAsBytes, err := hex.DecodeString(sW.PrivateKey)
	if err != nil {
		Error.Fatal(err.Error())
		os.Exit(1)
	}

	w.PrivateKey.D = new(big.Int).SetBytes(privateKeyAsBytes)
	w.PrivateKey.PublicKey.Curve = curve
	w.PrivateKey.PublicKey.X, w.PrivateKey.PublicKey.Y = curve.ScalarBaseMult(privateKeyAsBytes)
	w.PublicKey, err = hex.DecodeString(sW.PublicKey)
	if err != nil {
		Error.Fatal(err.Error())
		os.Exit(1)
	}

	w.Address = sW.Address
	return w
}

func (sW StorableWallet) String() string {
	strWallet := fmt.Sprint("\n  ** Wallet Information ** \n")
	strWallet += fmt.Sprintf("  + Private Key (%d bytes) : %s\n", len(sW.PrivateKey), sW.PrivateKey)
	strWallet += fmt.Sprintf("  + Public Key (%d bytes) : %s\n", len(sW.PublicKey), sW.PublicKey)
	strWallet += fmt.Sprintf("  + Address (%d bytes) : %s\n", len(sW.Address), sW.Address)
	return strWallet
}

func (w Wallet) String() string {
	bs, err := json.MarshalIndent(w, "", "   ")
	if err != nil {
		Error.Fatal(err.Error())
		os.Exit(1)
	}
	return string(bs)
}
