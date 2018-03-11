package main

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

// UTXOSet represents UTXO set
type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string]TxOutputMap) {
	unspentOutputs := make(map[string]TxOutputMap)
	var accumulated int
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			txOutputMap := DeserializeTxOutputMap(v)

			for outIdx, out := range txOutputMap {
				if out.isLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					if unspentOutputs[txID] == nil {
						unspentOutputs[txID] = make(TxOutputMap)
					}
					unspentOutputs[txID][outIdx] = out
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) TxOutputMap {
	UTXOs := make(TxOutputMap)
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeTxOutputMap(v)

			for idx, out := range outs {
				if out.isLockedWithKey(pubKeyHash) {
					UTXOs[idx] = out
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.findUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

func (u UTXOSet) getAllAddressInfo() map[string]int {
	db := u.Blockchain.db
	UTOX := make(map[string]TxOutputMap)
	addressInfos := make(map[string]int)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			txOutputs := DeserializeTxOutputMap(v)
			UTOX[txID] = txOutputs
		}
		return nil
	})

	for _, txOutputs := range UTOX {
		for _, txOutput := range txOutputs {
			address := hex.EncodeToString(txOutput.PubKeyHash)

			addressInfos[address] += txOutput.Value
		}
	}

	return addressInfos
}

func (u UTXOSet) getTotalValueOwnBy(publicKeyHash []byte) int {
	db := u.Blockchain.db
	totalValue := 0

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txOutputs := DeserializeTxOutputMap(v)
			for _, output := range txOutputs {
				if bytes.Compare(output.PubKeyHash, publicKeyHash) == 0 {
					totalValue += output.Value
				}
			}
		}
		return nil
	})

	return totalValue
}

func (u UTXOSet) verifyTxInputs(txIns []TxInput) bool {
	db := u.Blockchain.db
	isValid := true

	db.View(func(boltTx *bolt.Tx) error {
		b := boltTx.Bucket([]byte(utxoBucket))

		for _, vin := range txIns {
			outsBytes := b.Get(vin.Txid)

			if outsBytes != nil {
				outs := DeserializeTxOutputMap(outsBytes)
				if _, ok := outs[vin.TxOutIdx]; !ok {
					isValid = false
					return nil
				}
			} else {
				isValid = false
				return nil
			}
		}

		return nil
	})

	return isValid
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if tx.isCoinbase() == false {
				for _, vin := range tx.TxIns {
					updatedOuts := make(TxOutputMap)
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeTxOutputMap(outsBytes)

					for outIdx, out := range outs {
						if outIdx != vin.TxOutIdx {
							updatedOuts[outIdx] = out
						}
					}

					if len(updatedOuts) == 0 {
						err := b.Delete(vin.Txid)
						if err != nil {
							log.Panic(err)
						}
					} else {
						err := b.Put(vin.Txid, updatedOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}

				}
			}

			newOutputs := make(TxOutputMap)
			for idx, out := range tx.TxOuts {
				newOutputs[idx] = out
			}

			err := b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
