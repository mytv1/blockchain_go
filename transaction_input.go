package main

import (
	"fmt"
)

// TxInput is input component of a transaction
type TxInput struct {
	Txid      []byte `json:"Txid"`
	TxOutIdx  int    `json:"TxOutIdx"`
	Signature []byte `json:"Signature"`
}

func (txInput TxInput) String() string {
	str := fmt.Sprintf("Txid : %x\n", txInput.Txid)
	str += fmt.Sprintf("      TxOutIdx : %d\n", txInput.TxOutIdx)
	str += fmt.Sprintf("      Signature : %x\n", txInput.Signature)
	return str
}

// func (in *TxInput) usesKey(publicKeyHash []byte) bool {
// 	lockingHash := hashPublicKey(in.PubKey)
// 	return bytes.Compare(lockingHash, publicKeyHash) == 0
// }
