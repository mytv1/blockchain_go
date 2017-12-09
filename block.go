package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Hash          []byte
	PrevBlockHash []byte
	Data          []byte
	Timestamp     int64
}

func (b *Block) SetHash() {
	bTimeStamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	blockAsBytes := bytes.Join([][]byte{b.PrevBlockHash, b.Data, bTimeStamp}, []byte{})
	hash := sha256.Sum256(blockAsBytes)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{[]byte{}, prevBlockHash, []byte(data), time.Now().Unix()}
	block.SetHash()
	return block
}

func NewGenesisBlock(starting string) *Block {
	return NewBlock(starting, []byte{})
}
