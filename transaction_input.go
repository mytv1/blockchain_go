package main

import (
	"fmt"
)

// TxInput is input component of a transaction
type TxInput struct {
	Txid      []byte `json:"Txid"`
	TxOutIdx  int    `json:"TxOutIdx"`
	Signature []byte `json:"Signature"`
	PubKey    []byte `json:"Pubkey"`
}

func (txInput TxInput) String() string {
	str := fmt.Sprintf("Txid : %x\n", txInput.Txid)
	str += fmt.Sprintf("      TxOutIdx : %d\n", txInput.TxOutIdx)
	str += fmt.Sprintf("      Signature : %x\n", txInput.Signature)
	str += fmt.Sprintf("      Public key : %x\n", txInput.PubKey)
	return str
}

// func (in *TxInput) usesKey(publicKeyHash []byte) bool {
// 	lockingHash := hashPublicKey(in.PubKey)
// 	return bytes.Compare(lockingHash, publicKeyHash) == 0
// }
