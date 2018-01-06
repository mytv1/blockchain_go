package main

import (
	"encoding/json"
	"os"
)

const (
	// CmdSpreadHashList is used to spread hash list msg to other nodes
	CmdSpreadHashList = "SPR_HL"

	// CmdReqBlockchain is used to request entire blockchain
	CmdReqBlockchain = "REQ_BC"
	// CmdReqBestHeight is used to request node's current height
	CmdReqBestHeight = "REQ_BH"
	// CmdReqBlock is used to request a single block
	CmdReqBlock = "REQ_BL"
	// CmdPrintBlockchain is used to request node's blockchain printing
	CmdPrintBlockchain = "REQ_PRINT_BC"
	// CmdReqAddBlock is used to request node to add a block
	CmdReqAddBlock = "REQ_ADD_BL"

	// CmdResBestHeight is used to reponse with its own blockchain height
	CmdResBestHeight = "RES_BH"
	// CmdResBlock is used to response with node's single block
	CmdResBlock = "RES_BL"
	// CmdResBlockchain is used to reponse with entire blockchain
	CmdResBlockchain = "RES_BC"
)

// Message is used to communicate between nodes
type Message struct {
	Cmd    string `json:"Cmd"`
	Data   []byte `json:"Data"`
	Source Node   `json:"Source"`
}

func createMessageBc(bc *Blockchain) *Message {
	m := new(Message)
	m.Cmd = CmdResBlockchain
	m.Source = getLocalNode()
	m.Data = bc.serialize()
	return m
}

func createMsRequestBestHeight() *Message {
	m := new(Message)
	m.Source = getLocalNode()
	m.Cmd = CmdReqBestHeight
	return m
}

func createMsReponseBestHeight(bestHeight uint8) *Message {
	m := new(Message)
	m.Cmd = CmdResBestHeight
	m.Source = getLocalNode()
	m.Data = []byte{bestHeight}
	return m
}

func createMsRequestBlock(index uint8) *Message {
	m := new(Message)
	m.Cmd = CmdReqBlock
	m.Source = getLocalNode()
	m.Data = append(m.Data, byte(index))
	return m
}

func createMsResponseBlock(block *Block) *Message {
	m := new(Message)
	m.Cmd = CmdResBlock
	m.Source = getLocalNode()
	m.Data = block.serialize()
	return m
}

func createMsSpreadHashList(hashList [][]byte) *Message {
	m := new(Message)
	m.Cmd = CmdSpreadHashList
	m.Source = getLocalNode()
	data, err := json.Marshal(hashList)
	if err != nil {
		Error.Panic("Marshal fail")
		os.Exit(1)
	}
	m.Data = data
	return m
}

func (m *Message) serialize() []byte {
	data, err := json.Marshal(m)

	if err != nil {
		Error.Printf("Marshal fail")
		os.Exit(1)
	}
	return data
}

func deserializeMessage(data []byte) *Message {
	m := new(Message)
	err := json.Unmarshal(data, m)

	if err != nil {
		Error.Printf("Unmarshal fail")
		os.Exit(1)
	}

	return m
}
