package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const difficulty = 16

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func newProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	txAsBytes := []byte{}
	for _, tx := range pow.block.Transactions {
		txAsBytes = append(txAsBytes, tx.serialize()...)
	}

	data := bytes.Join(
		[][]byte{
			txAsBytes,
			pow.block.Header.PrevBlockHash,
			intToBytes(int(pow.block.Header.Timestamp)),
			intToBytes(pow.block.Header.Height),
			intToBytes(nonce),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	Info.Println("Mining...")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Header.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
