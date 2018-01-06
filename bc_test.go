package main

import "fmt"

func ExampleinitBlockchain() {
	bc := initBlockchain()
	bc.AddBlock("A send 1 btc to B")

	for _, block := range bc.blocks {
		fmt.Printf("Data : %s\n", block.Data)
	}

	// Output:
	// Data : Genesis block
	// Data : A send 1 btc to B
}
