package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Blockchain just array of blocks
type Blockchain struct {
	Blocks []*Block `json:"Blocks"`
}

var instantiated *Blockchain
var once sync.Once

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

func (bc *Blockchain) addBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) isEmpty() bool {
	return bc.Blocks == nil || bc.getBestHeight() == 0
}

func InitBlockchain() *Blockchain {
	once.Do(func() {
		instantiated = &Blockchain{[]*Block{newGenesisBlock("Genesis block")}}
	})
	return instantiated
}

func (bc *Blockchain) getBestHeight() int {
	return len(bc.Blocks)
}

func (bc *Blockchain) getHashList() [][]byte {
	var hashList [][]byte
	for _, block := range bc.Blocks {
		hashList = append(hashList, block.Hash)
	}
	return hashList
}

func (bc Blockchain) serialize() []byte {
	data, err := json.Marshal(bc)

	if err != nil {
		Error.Printf("Marshal fail\n")
		os.Exit(1)
	}

	return data
}

func deserialize(data []byte) *Blockchain {
	bc := new(Blockchain)
	err := json.Unmarshal(data, bc)

	if err != nil {
		Error.Printf("Unmarshal fail\n")
		os.Exit(1)
	}

	return bc
}
