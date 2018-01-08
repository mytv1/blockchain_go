package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

/*
Block simple structure
*/
type Block struct {
	Timestamp     int64  `json:"Timestamp"`
	Data          []byte `json:"Data"`
	PrevBlockHash []byte `json:"PrevBlockHash"`
	Hash          []byte `json:"Hash"`
	Height        int    `json:"Height"`
	Nonce         int    `json:"Nonce"`
}

func (b Block) String() string {
	var strBlock string
	strBlock += fmt.Sprintf("Prev hash: %x\n", b.PrevBlockHash)
	strBlock += fmt.Sprintf("Data: %s\n", b.Data)
	strBlock += fmt.Sprintf("Hash: %x\n", b.Hash)
	strBlock += fmt.Sprintf("Nonce: %x\n", b.Nonce)
	strBlock += fmt.Sprintf("Height: %x\n", b.Height)
	return strBlock
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// mine block
func newBlock(data string, prevBlockHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, height, 0}
	block.setHash()
	return block
}

func (b *Block) isGenesisBlock() bool {
	return len(b.PrevBlockHash) == 0
}

func newGenesisBlock() *Block {
	// add proof of work
	return newBlock("Genesis block", []byte{}, 1)
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
