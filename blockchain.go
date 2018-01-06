package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Blockchain just array of blocks
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

func (bc *Blockchain) addBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func newBlockchain() *Blockchain {
	return &Blockchain{[]*Block{newGenesisBlock()}}
}

func initBlockchain() {
	blockchain = newBlockchain()
}

func getBlockchain() *Blockchain {
	return blockchain
}

func setBlockchain(bc *Blockchain) {
	blockchain = bc
}

func (bc *Blockchain) getBestHeight() uint8 {
	return uint8(len(bc.Blocks))
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
