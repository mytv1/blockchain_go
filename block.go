package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

/*
Block simple structure
*/
type Block struct {
	Transactions []Transaction `json:"Transactions"`
	Header       Header        `json:"BlockHeader"`
}

/*Header of block */
type Header struct {
	Timestamp     int64  `json:"Timestamp"`
	Hash          []byte `json:"Hash"`
	PrevBlockHash []byte `json:"PrevBlockHash"`
	Height        int    `json:"Height"`
	Nonce         int    `json:"Nonce"`
}

func (b Block) String() string {
	var strBlock string
	strBlock += fmt.Sprintf("Prev hash: %x\n", b.Header.PrevBlockHash)
	strBlock += fmt.Sprintf("Transactions: \n")
	for idx, tx := range b.Transactions {
		strBlock += fmt.Sprintf("  Tx[%d] : %s\n", idx, tx)
	}
	strBlock += fmt.Sprintf("Hash: %x\n", b.Header.Hash)
	strBlock += fmt.Sprintf("Nonce: %d\n", b.Header.Nonce)
	strBlock += fmt.Sprintf("Height: %d\n", b.Header.Height)
	strBlock += fmt.Sprintf("Timestamp: %d\n", b.Header.Timestamp)
	return strBlock
}

func (b *Block) setHash() {
	hash := sha256.Sum256(b.serialize())

	b.Header.Hash = hash[:]
}

func newBlock(txs []Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{txs, Header{time.Now().Unix(), []byte{}, prevBlockHash, height, 0}}
	block.setHash()
	return block
}

func (b *Block) isGenesisBlock() bool {
	return len(b.Header.PrevBlockHash) == 0
}

func newGenesisBlock(txs []Transaction) *Block {
	return newBlock(txs, []byte{}, 1)
}

func (b *Block) hashTransactions() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(b.Transactions)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (b *Block) serialize() []byte {
	data, err := json.Marshal(b)

	if err != nil {
		Error.Printf("Marshal block fail\n")
		os.Exit(1)
	}
	return data
}

func deserializeBlock(data []byte) *Block {
	b := new(Block)
	err := json.Unmarshal(data, b)

	if err != nil {
		Error.Panic(err)
		os.Exit(1)
	}

	return b
}

func (h *Header) serialize() []byte {
	data, err := json.Marshal(h)

	if err != nil {
		Error.Printf("Marshal block fail\n")
		os.Exit(1)
	}
	return data
}

func deserializeHeader(data []byte) *Header {
	h := new(Header)
	err := json.Unmarshal(data, h)

	if err != nil {
		Error.Panic(err)
		os.Exit(1)
	}

	return h
}
