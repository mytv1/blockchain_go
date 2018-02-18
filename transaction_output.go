package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

type TxOutput struct {
	Value         int
	PublicKeyHash []byte
}

func (txOut *TxOutput) lock(address []byte) {
	decodedAddr := base58.Decode(address)
	publicKeyHash := decodedAddr[1 : len(decodedAddr)-4]
	out.PublicKeyHash = publicKeyHash
}

func (out TxOutput) isLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Compare(out.PublicKeyHash, publicKeyHash) == 0
}

func newTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

type TxOutputs struct {
	Outputs []TxOutput
}

func (outs TxOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
