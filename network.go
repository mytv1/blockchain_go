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

func getNeighborBc() {
	Info.Printf("Pull blockchain from other node in network...")
	network := getNetwork()

	for i := 0; i < maxAskTime; i++ {
		for _, node := range network.NeighborNodes {
			time.Sleep(1000 * time.Millisecond)
			if getBlockchain() == nil {
				sendRequestBc(node, nil)
			}
		}
	}

	if getBlockchain() == nil {
		Info.Printf("Pull failed. Create new blockchain.")
		initBlockchain()
	} else {
		bc := getBlockchain()
		Info.Printf("Pull completed. Blockchain height: %d", bc.getBestHeight())
	}
}

func sendRequestBc(node Node, bc *Blockchain) {
	var myHeight uint8
	if bc == nil || bc.getBestHeight() == 0 {
		bc = new(Blockchain)
		myHeight = 0
	} else {
		myHeight = uint8(bc.getBestHeight())
	}

	neighborHeight, err := getNeighborBcBestHeight(node)

	if err != nil {
		return
	}

	// Get blocks until local blockchain's height equal to neighbor's
	for myHeight < neighborHeight {
		ms := createMsRequestBlock(myHeight + uint8(1))
		data := ms.serialize()

		conn, err := net.Dial("tcp", node.Address)

		if err != nil {
			Error.Printf("%s is not avaiable\n", node.Address)
			return
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
		bc.Blocks = append(bc.Blocks, block)

		myHeight++
	}
	setBlockchain(bc)
}

func getNeighborBcBestHeight(node Node) (uint8, error) {
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
	scanner.Scan()
	msAsBytes := scanner.Bytes()

	message := deserializeMessage(msAsBytes)
	neighborHeigh := uint8(message.Data[0])
	return neighborHeigh, nil
}

func spreadHashList() {
	nw := getNetwork()
	bc := getBlockchain()

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
