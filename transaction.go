package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsidy = 25

// Transaction represent a transaction between wallets
type Transaction struct {
	ID     []byte     `json:"ID"`
	TxIns  []TxInput  `json:"TxIns"`
	TxOuts []TxOutput `json:"TxOuts"`
}

func (tx Transaction) isCoinbase() bool {
	return len(tx.TxIns) == 1 && len(tx.TxIns[0].Txid) == 0 && tx.TxIns[0].TxOutIdx == -1
}

func (tx Transaction) serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func deserializeTransaction(data []byte) *Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return &transaction
}

func (tx *Transaction) hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.serialize())

	return hash[:]
}

func (tx *Transaction) trimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, vin := range tx.TxIns {
		inputs = append(inputs, TxInput{vin.Txid, vin.TxOutIdx, nil, vin.PubKey})
	}

	for _, vout := range tx.TxOuts {
		outputs = append(outputs, TxOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

func (tx *Transaction) sign(privateKey ecdsa.PrivateKey) {
	if tx.isCoinbase() {
		return
	}

	txCopy := tx.trimmedCopy()

	dataToSign := fmt.Sprintf("%x", txCopy)
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
	signature := append(r.Bytes(), s.Bytes()...)

	if err != nil {
		Error.Panic(err)
	}

	for inID := range txCopy.TxIns {
		txCopy.TxIns[inID].Signature = nil
		tx.TxIns[inID].Signature = signature
	}
}

func (tx *Transaction) verifySignature() bool {
	txCopy := tx.trimmedCopy()
	curve := elliptic.P256()

	for _, vin := range tx.TxIns {
		publicKey := vin.PubKey

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(publicKey) - 1
		x.SetBytes(publicKey[1:(keyLen/2 + 1)])
		y.SetBytes(publicKey[(keyLen/2 + 1):])

		dataToVerify := fmt.Sprintf("%x", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
	}
	return true
}

func (tx *Transaction) verifyValues(prevTxs map[string]Transaction) bool {
	allInputValues := 0
	allOutputValues := 0
	for _, vin := range tx.TxIns {
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		allInputValues += prevTx.TxOuts[vin.TxOutIdx].Value
	}

	for _, vout := range tx.TxOuts {
		allOutputValues += vout.Value
	}

	return allInputValues == allOutputValues
}

func newCoinbaseTx(addrTo string) *Transaction {
	txIn := TxInput{[]byte{}, -1, nil, []byte{}}
	txOut := newTxOutput(subsidy, addrTo)
	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}
	tx.ID = tx.hash()

	return &tx
}

func (tx Transaction) String() string {
	strTx := fmt.Sprintf("\n    ID: %x\n", tx.ID)
	strTx += fmt.Sprintf("    Vin :\n")
	for idx, txIn := range tx.TxIns {
		strTx += fmt.Sprintf("      [%d]%v\n", idx, txIn)
	}

	strTx += fmt.Sprintf("    Vout :\n")
	for idx, txOut := range tx.TxOuts {
		strTx += fmt.Sprintf("      [%d]%v\n", idx, txOut)
	}

	return strTx
}
