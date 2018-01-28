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
	Data   []byte `json:"Data"`
	Header Header `json:"BlockHeader"`
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
	strBlock += fmt.Sprintf("Data: %s\n", b.Data)
	strBlock += fmt.Sprintf("Hash: %x\n", b.Header.Hash)
	strBlock += fmt.Sprintf("Nonce: %d\n", b.Header.Nonce)
	strBlock += fmt.Sprintf("Height: %d\n", b.Header.Height)
	strBlock += fmt.Sprintf("Timestamp: %d\n", b.Header.Timestamp)
	return strBlock
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Header.Timestamp, 10))
	headers := bytes.Join([][]byte{b.Header.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Header.Hash = hash[:]
}

// mine block
func newBlock(data string, prevBlockHash []byte, height int) *Block {
	block := &Block{[]byte(data), Header{time.Now().Unix(), []byte{}, prevBlockHash, height, 0}}
	block.setHash()
	return block
}

func (b *Block) isGenesisBlock() bool {
	return len(b.Header.PrevBlockHash) == 0
}

func newGenesisBlock() *Block {
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
