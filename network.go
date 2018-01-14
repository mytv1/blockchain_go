package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
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

func getNeighborBc() *Blockchain {
	var bc *Blockchain
	Info.Printf("Pull blockchain from other node in network...")
	network := getNetwork()

	for i := 0; i < maxAskTime; i++ {
		for _, node := range network.NeighborNodes {
			if bc == nil || bc.isEmpty() {
				bc = sendRequestBc(node, nil)
				if bc != nil && !bc.isEmpty() {
					Info.Printf("Pull completed. Blockchain height: %d", bc.getBestHeight())
					return bc
				}
			}
		}
	}

	return bc
}

func sendRequestBc(node Node, bc *Blockchain) *Blockchain {
	var myHeight int
	if bc == nil || bc.getBestHeight() == 0 {
		bc = new(Blockchain)
		myHeight = 0
	} else {
		myHeight = bc.getBestHeight()
	}

	neighborHeight, err := getNeighborBcBestHeight(node)

	if err != nil {
		return nil
	}

	// Get blocks until local blockchain's height equal to neighbor's
	for myHeight < neighborHeight {
		ms := createMsRequestBlock(myHeight + 1)
		data := ms.serialize()

		conn, err := net.Dial("tcp", node.Address)

		if err != nil {
			Error.Printf("%s is not avaiable\n", node.Address)
			return nil
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

	return bc
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
	scanner.Scan()
	msAsBytes := scanner.Bytes()

	message := deserializeMessage(msAsBytes)
	neighborHeight, err := strconv.Atoi(string(message.Data))
	if err != nil {
		Error.Printf("")
	}

	return neighborHeight, nil
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
