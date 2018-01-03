package main

import (
	"encoding/json"
	"os"
)

const (
	CMD_SPREAD_HASHLIST = "SPR_HL"

	CMD_REQ_BLOCKCHAIN   = "REQ_BC"
	CMD_REQ_BEST_HEIGHT  = "REQ_BH"
	CMD_REQ_BLOCK        = "REQ_BL"
	CMD_PRINT_BLOCKCHAIN = "REQ_PRINT_BC"
	CMD_REQ_ADD_BLOCK    = "REQ_ADD_BL"

	CMD_RES_BEST_HEIGHT = "RES_BH"
	CMD_RES_BLOCK       = "RES_BL"
	CMD_RES_BLOCKCHAIN  = "RES_BC"
)

type Message struct {
	Cmd    string `json:"Cmd"`
	Data   []byte `json:"Data"`
	Source Node   `json:"Source"`
}

func CreateMessageBc(bc *Blockchain) *Message {
	var m *Message = new(Message)
	m.Cmd = CMD_RES_BLOCKCHAIN
	m.Source = GetLocalNode()
	m.Data = bc.SerializeBlockchain()
	return m
}

func CreateMsRequestBestHeight() *Message {
	var m *Message = new(Message)
	m.Source = GetLocalNode()
	m.Cmd = CMD_REQ_BEST_HEIGHT
	return m
}

func CreateMsReponseBestHeight(bestHeight uint8) *Message {
	var m *Message = new(Message)
	m.Cmd = CMD_RES_BEST_HEIGHT
	m.Source = GetLocalNode()
	m.Data = []byte{bestHeight}
	return m
}

func CreateMsRequestBlock(index uint8) *Message {
	var m *Message = new(Message)
	m.Cmd = CMD_REQ_BLOCK
	m.Source = GetLocalNode()
	m.Data = append(m.Data, byte(index))
	return m
}

func CreateMsResponseBlock(block *Block) *Message {
	var m *Message = new(Message)
	m.Cmd = CMD_RES_BLOCK
	m.Source = GetLocalNode()
	m.Data = block.Serialize()
	return m
}

func CreateMsSpreadHashList(hashList [][]byte) *Message {
	var m *Message = new(Message)
	m.Cmd = CMD_SPREAD_HASHLIST
	m.Source = GetLocalNode()
	data, err := json.Marshal(hashList)
	if err != nil {
		Error.Panic("Marshal fail")
		os.Exit(1)
	}
	m.Data = data
	return m
}

func (m *Message) Serialize() []byte {
	data, err := json.Marshal(m)

	if err != nil {
		Error.Printf("Marshal fail")
		os.Exit(1)
	}
	return data
}

func DeserializeMessage(data []byte) *Message {
	var m *Message = new(Message)
	err := json.Unmarshal(data, m)

	if err != nil {
		Error.Printf("Unmarshal fail")
		os.Exit(1)
	}

	return m
}
