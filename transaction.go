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

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

func (tx Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx Transaction) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) hash() byte {
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
		inputs = append(inputs, TxInput{vin.Txid, vin.Vout, nil, vin.PubicKey})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TxOutput{vout.Value, vout.PublicKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

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
	signature := append(r.Bytes(), s.Bytes()...)

	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
	if err != nil {
		Error.Panic(err)
	}

	for inID, vin := range txCopy.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubicKey = nil
	}
}
