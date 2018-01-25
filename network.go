package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"strconv"
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

func syncWithNeighborNode(bc *Blockchain) {
	Info.Printf("Pull blockchain from other node in network...")
	network := getNetwork()

	for i := 0; i < maxAskTime; i++ {
		for _, node := range network.NeighborNodes {
			if sendRequestBc(node, bc) {
				Info.Printf("Pull completed. Blockchain height: %d", bc.getBestHeight())
				return
			}
		}
	}
}

func sendRequestBc(node Node, bc *Blockchain) bool {
	myHeight := bc.getBestHeight()

	neighborHeight, err := getNeighborBcBestHeight(node)

	if err != nil {
		return false
	}

	Info.Printf(" %d - %d ", myHeight, neighborHeight)
	minHeight := min(myHeight, neighborHeight)

	for i := 1; i <= minHeight; i++ {
		if compareBlockWithNeighbor(bc.getBlockByHeight(i), node) {
			// TODO write sth
			break
		} else {
			Error.Fatal("Independent blockchain detected. Program exit.")
			os.Exit(1)
		}
	}

	if myHeight < neighborHeight {
		for i := myHeight + 1; i <= neighborHeight; i++ {
			pullBlockFromNeighbor(bc, node, i)
		}
	} else if myHeight > neighborHeight {
		// TODO spread height
	}
	return true
}

func compareBlockWithNeighbor(b *Block, node Node) bool {
	ms := createMsRequestHeaderValidation(b.Header)
	data := ms.serialize()

	conn, err := net.Dial("tcp", node.Address)

	if err != nil {
		Error.Printf("%s is not avaiable\n", node.Address)
		return false
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

	Error.Printf("Message : %s", msAsBytes)
	isValid, err := strconv.ParseBool(string(msResponse.Data))

	if err != nil {
		Error.Printf("Parse failed")
		return false
	}
	return isValid
}

func pullBlockFromNeighbor(bc *Blockchain, node Node, blockIndex int) {
	ms := createMsRequestBlock(blockIndex)
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
	bc.addBlock(block)
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
