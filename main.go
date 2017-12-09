package main

import "fmt"

func main() {
	bc := NewBlockChain("Genesis block")

	bc.AddBlock("Doraemon send 1 btc to batman")
	bc.AddBlock("Batman send 2 btc to superman")
	bc.AddBlock("Batman send 1 btc to girls")

	fmt.Println()
	for _, block := range bc.blocks {
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Timestamp : %x\n", block.Timestamp)
		fmt.Printf("Previous Hash : %x\n", block.PrevBlockHash)
		fmt.Println()
	}
}
