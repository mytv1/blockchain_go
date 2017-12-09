package main

type Blockchain struct {
	blocks []*Block
}

func NewBlockChain(starting string) *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock(starting)}}
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
