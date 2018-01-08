package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

const dbFileName = "bc.db"
const blocksBucketName = "blocks"

// Blockchain implement interactions with a DB
type Blockchain struct {
	db *bolt.DB
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	bc          *Blockchain
}

func (i *BlockchainIterator) next() *Block {
	var block *Block

	err := i.bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucketName))
		encodedBlock := b.Get(i.currentHash)
		block = deserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		Error.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

func (bc *Blockchain) iterator() *BlockchainIterator {
	var lastHash = bc.getTopBlockHash()
	bci := &BlockchainIterator{lastHash, bc}
	return bci
}

func (bc *Blockchain) getTopBlockHash() []byte {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucketName))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		Error.Panic(err)
	}

	return lastHash
}

func (bc *Blockchain) String() string {
	bci := bc.iterator()
	var strBlockchain string

	for {
		block := bci.next()
		strBlock := fmt.Sprintf("%v", block)
		strBlockchain += "[" + strconv.Itoa(block.Height) + "]  "
		strBlockchain += strBlock
		strBlockchain += "\n"

		if block.isGenesisBlock() {
			break
		}
	}

	return strBlockchain
}

func (bc *Blockchain) addBlock(block *Block) {
	pow := newProofOfWork(block)
	// os.Exit(1)
	nonce, hash := pow.run()
	block.Nonce = nonce
	block.Hash = hash[:]

	Info.Printf(" %v ", block)

	if bc.isBlockchainEmpty() {
		bc.addGenesisBlock(block)
	} else {
		bc.addNormalBlock(block)
	}
}

func (bc *Blockchain) addGenesisBlock(genesisBlock *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucketName))

		err := b.Put(genesisBlock.Hash, genesisBlock.serialize())
		if err != nil {
			Error.Panic(err)
		}

		err = b.Put([]byte("l"), genesisBlock.Hash)
		if err != nil {
			Error.Panic(err)
		}

		return nil
	})

	if err != nil {
		Error.Panic(err)
	}
}

func (bc *Blockchain) addNormalBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucketName))

		lastHash := b.Get([]byte("l"))
		encodedLastBlock := b.Get(lastHash)
		lastBlock := deserializeBlock(encodedLastBlock)

		// TODO : Add verify function
		if block.Height > lastBlock.Height && bytes.Compare(block.PrevBlockHash, lastBlock.Hash) == 0 {
			blockData := block.serialize()
			err := b.Put(block.Hash, blockData)
			if err != nil {
				Error.Panic(err)
			}

			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				Error.Panic(err)
			}
		}
		return nil
	})

	if err != nil {
		Error.Panic(err)
	}
}

func createEmptyBlockchain() *Blockchain {
	if isBLockchainExist() {
		fmt.Println("Blockchain already exists.")
		return nil
	}

	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		Error.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(blocksBucketName))
		if err != nil {
			Error.Panic(err)
		}

		return nil
	})

	if err != nil {
		Error.Fatal(err)
	}

	bc := &Blockchain{db}
	return bc
}

func (bc *Blockchain) isBlockchainEmpty() bool {
	return len(bc.getTopBlockHash()) == 0
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) getBestHeight() int {
	var lastBlock *Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucketName))
		lastHash := b.Get([]byte("l"))
		if lastHash == nil {
			return nil
		}

		blockData := b.Get(lastHash)
		lastBlock = deserializeBlock(blockData)

		return nil
	})

	if err != nil {
		Error.Panic(err)
		return 0
	}

	if lastBlock == nil {
		return 0
	}

	return lastBlock.Height
}

func (bc *Blockchain) getHashList() [][]byte {
	var hashList [][]byte
	bci := bc.iterator()

	for {
		block := bci.next()

		hashList = append(hashList, block.Hash)

		if block.isGenesisBlock() {
			break
		}
	}
	return hashList
}

func isBLockchainExist() bool {
	return isDbExists(dbFileName)
}

func isDbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (bc *Blockchain) getBlockByHeight(height int) *Block {
	bci := bc.iterator()

	for {
		block := bci.next()

		if block.Height == height {
			return block
		}

		if block.isGenesisBlock() {
			break
		}
	}

	return nil
}

func getLocalBc() *Blockchain {
	if !isDbExists(dbFileName) {
		Info.Printf("Local blockchain not exists")
		return nil
	}

	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		Error.Fatal(err)
	}

	bc := &Blockchain{db}
	return bc
}
