package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
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

	i.currentHash = block.Header.PrevBlockHash

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
		strBlockchain += "[" + strconv.Itoa(block.Header.Height) + "]  "
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

	if !pow.validate() {
		nonce, hash := pow.run()
		block.Header.Nonce = nonce
		block.Header.Hash = hash[:]
	}

	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucketName))
		if bc.isEmpty() {
			bc.putBlock(b, block.Header.Hash, block.serialize())
		} else {
			lastHash := b.Get([]byte("l"))
			encodedLastBlock := b.Get(lastHash)
			lastBlock := deserializeBlock(encodedLastBlock)

			if block.Header.Height > lastBlock.Header.Height && bytes.Compare(block.Header.PrevBlockHash, lastBlock.Header.Hash) == 0 {
				bc.putBlock(b, block.Header.Hash, block.serialize())
			} else {
				Error.Printf("Block invalid. Failed to add block. : \n%v\n", block)
				Error.Printf("Last bl. : \n%v\n", lastBlock)
			}
		}
		return nil
	})
	if err != nil {
		Error.Panic(err)
	}
}

func (bc *Blockchain) putBlock(b *bolt.Bucket, blockHash, blockData []byte) {
	err := b.Put(blockHash, blockData)
	if err != nil {
		Error.Panic(err)
	}

	err = b.Put([]byte("l"), blockHash)
	if err != nil {
		Error.Panic(err)
	}
}

func createEmptyBlockchain() *Blockchain {
	if isDbExists(dbFileName) {
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

func (bc *Blockchain) isEmpty() bool {
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

	return lastBlock.Header.Height
}

func (bc *Blockchain) getHashList() [][]byte {
	var hashList [][]byte
	bci := bc.iterator()

	for {
		block := bci.next()

		hashList = append(hashList, block.Header.Hash)

		if block.isGenesisBlock() {
			break
		}
	}
	return hashList
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

		if block.Header.Height == height {
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
		return nil
	}

	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		Error.Fatal(err)
	}

	bc := &Blockchain{db}
	return bc
}

func (bc *Blockchain) findUTXO() map[string]TxOutputMap {
	UTXO := make(map[string]TxOutputMap)
	spentTXOs := make(map[string][]int)
	bci := bc.iterator()

	for {
		block := bci.next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.TxOuts {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				txOutputMap := UTXO[txID]
				if txOutputMap == nil {
					txOutputMap = make(TxOutputMap)
				}
				txOutputMap[outIdx] = out
				UTXO[txID] = txOutputMap
			}

			if tx.isCoinbase() == false {
				for _, in := range tx.TxIns {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.TxOutIdx)
				}
			}
		}

		if len(block.Header.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

func (bc *Blockchain) newTransaction(wallet *Wallet, to string, amount int) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()
	pubKeyHash := hashPublicKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for idx := range outs {
			input := TxInput{txID, idx, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", wallet.Address)
	outputs = append(outputs, *newTxOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *newTxOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.hash()
	tx.sign(wallet.PrivateKey)

	return &tx
}

func (bc *Blockchain) verifyTransaction(tx *Transaction) bool {
	if tx.isCoinbase() {
		return true
	}

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	prevTxs := bc.findTransactionsByTx(tx)

	return tx.verifySignature() && UTXOSet.verifyTxInputs(tx.TxIns) && tx.verifyValues(prevTxs)
}

func (bc *Blockchain) findTransactionsByTx(tx *Transaction) map[string]Transaction {
	prevTxs := make(map[string]Transaction)
	for _, vin := range tx.TxIns {
		prevTx, err := bc.findTransaction(vin.Txid)
		if err != nil {
			Error.Fatal(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return prevTxs
}

func (bc *Blockchain) findTransaction(ID []byte) (Transaction, error) {
	bci := bc.iterator()

	for {
		block := bci.next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return tx, nil
			}
		}

		if len(block.Header.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}
