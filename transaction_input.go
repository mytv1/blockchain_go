package main

import "bytes"

type TxInput struct {
	Txid       []byte
	TxOutIndex int
	Signature  []byte
}

func (in *TXInput) usesKey(publicKeyHash []byte) bool {
	lockingHash := hashPublicKey(in.PubKey)
	return bytes.Compare(lockingHash, publicKeyHash) == 0
}
