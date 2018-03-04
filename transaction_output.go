package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

// TxOutput is output component of a transaction
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

func (txOut *TxOutput) lock(address string) {
	decodedAddr := base58.Decode(address)
	pubKeyHash := decodedAddr[1 : len(decodedAddr)-4]
	txOut.PubKeyHash = pubKeyHash
}

func (txOut TxOutput) isLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txOut.PubKeyHash, pubKeyHash) == 0
}

func newTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.lock(address)
	return txo
}

func (txOut TxOutput) String() string {
	str := fmt.Sprintf("Value : %d\n", txOut.Value)
	str += fmt.Sprintf("      PubKeyHash : %x ", txOut.PubKeyHash)
	return str
}

// TxOutputs is a set of txoutput
type TxOutputs struct {
	Outputs []TxOutput
}

// Serialize serialize txoutouts
func (outs TxOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserialize txoutouts
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
