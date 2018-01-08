package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"
)

// Node contains its own address
type Node struct {
	Address string `json:"address"`
}

// Network contains runn other neighbor nodes
type Network struct {
	// LocalNode is running node by program
	LocalNode Node `json:"local_node"`

	// NeighborNodes are other nodes specified in config.json
	NeighborNodes []Node `json:"neighbor_nodes"`
}

const (
	maxAskTime = 2
)

func getNetwork() Network {
	return getConfig().Nw
}

func getLocalNode() Node {
	return getConfig().Nw.LocalNode
}

func getNeighborBc() *Blockchain {
	Info.Printf("Pull blockchain from other node in network...")
	network := getNetwork()
	bc := createEmptyBlockchain()

	for i := 0; i < maxAskTime; i++ {
		for _, node := range network.NeighborNodes {
			time.Sleep(1000 * time.Millisecond)
			if bc.isBlockchainEmpty() {
				bc = sendRequestBc(node, *bc)
			}
		}
	}

	if bc.isBlockchainEmpty() {
		bc.addBlock(newGenesisBlock())
	} else {
		Info.Printf("Pull completed. Blockchain height: %d", bc.getBestHeight())
	}

	return bc
}

func sendRequestBc(node Node, bc Blockchain) *Blockchain {
	myHeight := bc.getBestHeight()

	neighborHeight, err := getNeighborBcBestHeight(node)

	if err != nil {
		return &bc
	}

	// Get blocks until local blockchain's height equal to neighbor's
	for myHeight < neighborHeight {
		ms := createMsRequestBlock(myHeight + 1)
		data := ms.serialize()

		conn, err := net.Dial("tcp", node.Address)

		if err != nil {
			Error.Printf("%s is not avaiable\n", node.Address)
			return &bc
		}
		defer conn.Close()

		_, err = io.Copy(conn, bytes.NewReader(data))
		if err != nil {
			Error.Panic(err)
		}

		scanner := bufio.NewScanner(bufio.NewReader(conn))
		scanner.Scan()
		msAsBytes := scanner.Bytes()
		msResponse := deserializeMessage(msAsBytes)

		block := deserializeBlock(msResponse.Data)
		bc.addBlock(block)

		myHeight++
	}
	return &bc
}

func getNeighborBcBestHeight(node Node) (int, error) {
	m := createMsRequestBestHeight()

	conn, err := net.Dial("tcp", node.Address)

	if err != nil {
		Error.Printf("%s is not avaiable\n", node.Address)
		return 0, err
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(m.serialize()))
	if err != nil {
		Error.Panic(err)
		return 0, err
	}

	scanner := bufio.NewScanner(bufio.NewReader(conn))
	ok := scanner.Scan()
	if !ok {
		return 0, nil
	}

	msAsBytes := scanner.Bytes()

	message := deserializeMessage(msAsBytes)
	neighborHeigh := bytesToInt(message.Data)
	return neighborHeigh, nil
}

func spreadHashList(bc *Blockchain) {
	nw := getNetwork()

	for _, node := range nw.NeighborNodes {
		m := createMsSpreadHashList(bc.getHashList())
		sendMessage(node, m)
	}
}

func sendMessage(node Node, m *Message) {
	conn, err := net.Dial("tcp", node.Address)

	if err != nil {
		Error.Printf("%s is not avaiable\n", node.Address)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(m.serialize()))
	if err != nil {
		Error.Panic(err)
		return
	}
}
