package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Blockchain struct {
	Blocks []*Block `json:"Blocks"`
}

func (bc *Blockchain) String() string {
	var strBlockchain string
	for index, block := range bc.Blocks {
		strBlock := fmt.Sprintf("%v", block)
		strBlockchain += "[" + strconv.Itoa(index) + "]  "
		strBlockchain += strBlock
		strBlockchain += "\n"
	}
	return strBlockchain
}

var blockchain *Blockchain

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func InitBlockchain() {
	blockchain = NewBlockchain()
}

func GetBlockchain() *Blockchain {
	return blockchain
}

func SetBlockchain(bc *Blockchain) {
	blockchain = bc
}

func (b *Blockchain) GetBestHeight() uint8 {
	return uint8(len(b.Blocks))
}

func (bc *Blockchain) GetHashList() [][]byte {
	var hashList [][]byte
	for _, block := range bc.Blocks {
		hashList = append(hashList, block.Hash)
	}
	return hashList
}

func (bc Blockchain) SerializeBlockchain() []byte {
	data, err := json.Marshal(bc)

	if err != nil {
		Error.Printf("Marshal fail\n")
		os.Exit(1)
	}

	return data
}

func DeserializeBlockchain(data []byte) *Blockchain {
	var bc *Blockchain = new(Blockchain)
	err := json.Unmarshal(data, bc)

	if err != nil {
		Error.Printf("Unmarshal fail\n")
		os.Exit(1)
	}

	return bc
}
