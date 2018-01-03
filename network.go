package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"
)

type Node struct {
	Address string `json:"address"`
}

type Network struct {
	LocalNode     Node   `json:"local_node"`
	NeighborNodes []Node `json:"neighbor_nodes"`
}

const (
	MAX_ASK_TIME = 2
)

func GetNetwork() Network {
	return GetConfig().Nw
}

func GetLocalNode() Node {
	return GetConfig().Nw.LocalNode
}

func GetNeighborBc() {
	Info.Printf("Pull blockchain from other node in network...")
	network := GetNetwork()

	for i := 0; i < MAX_ASK_TIME; i++ {
		for _, node := range network.NeighborNodes {
			time.Sleep(1000 * time.Millisecond)
			if GetBlockchain() == nil {
				SendRequestBc(node, nil)
			}
		}
	}

	if GetBlockchain() == nil {
		Info.Printf("Pull failed. Create new blockchain.")
		InitBlockchain()
	} else {
		bc := GetBlockchain()
		Info.Printf("Pull completed. Blockchain height: %d", bc.GetBestHeight())
	}
}

func SendRequestBc(node Node, bc *Blockchain) {
	var myHeight uint8
	if bc == nil || bc.GetBestHeight() == 0 {
		bc = new(Blockchain)
		myHeight = 0
	} else {
		myHeight = uint8(bc.GetBestHeight())
	}

	neighborHeight, err := GetNeighborBcBestHeight(node)

	if err != nil {
		return
	}

	// Get blocks until local blockchain's height equal to neighbor's
	for myHeight < neighborHeight {
		ms := CreateMsRequestBlock(myHeight + uint8(1))
		data := ms.Serialize()

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
		msResponse := DeserializeMessage(msAsBytes)

		block := DeserializeBlock(msResponse.Data)
		bc.Blocks = append(bc.Blocks, block)

		myHeight++
	}
	SetBlockchain(bc)
}

func GetNeighborBcBestHeight(node Node) (uint8, error) {
	m := CreateMsRequestBestHeight()

	conn, err := net.Dial("tcp", node.Address)

	if err != nil {
		Error.Printf("%s is not avaiable\n", node.Address)
		return 0, err
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(m.Serialize()))
	if err != nil {
		Error.Panic(err)
		return 0, err
	}

	scanner := bufio.NewScanner(bufio.NewReader(conn))
	scanner.Scan()
	msAsBytes := scanner.Bytes()

	message := DeserializeMessage(msAsBytes)
	neighborHeigh := uint8(message.Data[0])
	return neighborHeigh, nil
}

func SpreadHashList() {
	nw := GetNetwork()
	bc := GetBlockchain()

	for _, node := range nw.NeighborNodes {
		m := CreateMsSpreadHashList(bc.GetHashList())
		SendMessage(node, m)
	}
}

func SendMessage(node Node, m *Message) {
	conn, err := net.Dial("tcp", node.Address)

	if err != nil {
		Error.Printf("%s is not avaiable\n", node.Address)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(m.Serialize()))
	if err != nil {
		Error.Panic(err)
		return
	}
}
