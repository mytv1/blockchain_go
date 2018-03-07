package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 25

// Transaction represent a transaction between wallets
type Transaction struct {
	ID   []byte     `json:"ID"`
	Vin  []TxInput  `json:"Vin"`
	Vout []TxOutput `json:"Vout"`
}

func (tx Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].TxOutIdx == -1
}

// Serialize serialize a transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *Transaction) trimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TxInput{vin.Txid, vin.TxOutIdx, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TxOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Sign make signature of a transaction
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.isCoinbase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTxs[hex.EncodeToString(vin.Txid)].ID == nil {
			Error.Panic("Previous transaction is not correct")
		}
	}

	txCopy := tx.trimmedCopy()

	dataToSign := fmt.Sprintf("%x", txCopy)
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
	signature := append(r.Bytes(), s.Bytes()...)

	if err != nil {
		Error.Panic(err)
	}

	for inID := range txCopy.Vin {
		txCopy.Vin[inID].Signature = nil
		tx.Vin[inID].Signature = signature
	}
}

func newCoinbaseTx(addrTo string) *Transaction {
	txIn := TxInput{[]byte{}, -1, nil}
	txOut := newTxOutput(subsidy, addrTo)
	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}
	tx.ID = tx.hash()

	return &tx
}

func (tx Transaction) String() string {
	strTx := fmt.Sprintf("\n    ID: %x\n", tx.ID)
	strTx += fmt.Sprintf("    Vin :\n")
	for idx, txIn := range tx.Vin {
		strTx += fmt.Sprintf("      [%d]%v", idx, txIn)
	}

	strTx += fmt.Sprintf("    Vout :\n")
	for idx, txOut := range tx.Vout {
		strTx += fmt.Sprintf("      [%d]%v", idx, txOut)
	}

	return strTx
}
