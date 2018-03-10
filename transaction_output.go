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
	Value      int    `json:"Value"`
	PubKeyHash []byte `json:"PubKeyHash"`
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

// TxOutputMap is map of TxOutput
type TxOutputMap map[int]TxOutput

// Serialize serialize txoutouts
func (txOutputMap *TxOutputMap) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(txOutputMap)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeTxOutputMap deserialize txoutouts
func DeserializeTxOutputMap(data []byte) TxOutputMap {
	var txOutputMap TxOutputMap

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&txOutputMap)
	if err != nil {
		log.Panic(err)
	}

	return txOutputMap
}
