package main

import "sync"

type Blockchain struct {
	blocks []*Block
}

var instantiated *Blockchain
var once sync.Once

func InitBlockchain() *Blockchain {
	once.Do(func() {
		instantiated = &Blockchain{[]*Block{NewGenesisBlock("Genesis block")}}
	})
	return instantiated
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
