package main

import (
	"encoding/json"
	"os"
	"strconv"
)

const (
	// PING

	// CmdSpreadHashList is used to spread hash list msg to other nodes
	CmdSpreadHashList = "SPR_HL"

	// CmdReqBestHeight is used to request node's current height
	CmdReqBestHeight = "REQ_BH"
	// CmdReqBlock is used to request a single block
	CmdReqBlock = "REQ_BL"
	// CmdPrintBlockchain is used to request node's blockchain printing
	CmdPrintBlockchain = "REQ_PRINT_BC"
	// CmdReqAddBlock is used to request node to add a block
	CmdReqAddBlock = "REQ_ADD_BL"
	// CmdReqAddress is used to request nodes address
	CmdReqAddress = "REQ_ADDR"
	// CmdReqHeaderValidation is used to request other node to valid a blocks hash list
	CmdReqHeaderValidation = "REQ_BL_VAL"
	// CmdReqPrintAllAddressInfo is used to request other node to valid a blocks hash list
	CmdReqPrintAllAddressInfo = "REQ_ALL_ADDR"

	// CmdResBestHeight is used to reponse with its own blockchain height
	CmdResBestHeight = "RES_BH"
	// CmdResBlock is used to response with node's single block
	CmdResBlock = "RES_BL"
	// CmdResHeaderValidation is used to response CmdReqHeaderValidation
	CmdResHeaderValidation = "RES_BL_VAL"
	// CmdResAddress is used to response CmdReqHeaderValidation
	CmdResAddress = "RES_ADDR"
)

// Message is used to communicate between nodes
type Message struct {
	Cmd    string `json:"Cmd"`
	Data   []byte `json:"Data"`
	Source Node   `json:"Source"`
}

func createBaseMessage() *Message {
	m := new(Message)
	m.Source = getLocalNode()
	return m
}

func createMsRequestBestHeight() *Message {
	m := createBaseMessage()
	m.Cmd = CmdReqBestHeight
	return m
}

func createMsReponseBestHeight(bestHeight int) *Message {
	m := createBaseMessage()
	m.Cmd = CmdResBestHeight
	m.Data = intToBytes(bestHeight)
	return m
}

func createMsRequestBlock(index int) *Message {
	m := createBaseMessage()
	m.Cmd = CmdReqBlock
	m.Data = intToBytes(index)
	return m
}

func createMsResponseBlock(block *Block) *Message {
	m := createBaseMessage()
	m.Cmd = CmdResBlock
	m.Data = block.serialize()
	return m
}

func createMsRequestAddress() *Message {
	m := createBaseMessage()
	m.Cmd = CmdReqAddress
	return m
}

func createMsResponseAddress() *Message {
	m := createBaseMessage()
	m.Cmd = CmdResAddress
	return m
}

func createMsSpreadHashList(hashList [][]byte) *Message {
	m := createBaseMessage()
	m.Cmd = CmdSpreadHashList
	data, err := json.Marshal(hashList)
	if err != nil {
		Error.Panic("Marshal fail")
		os.Exit(1)
	}
	m.Data = data
	return m
}

func createMsRequestHeaderValidation(blockHeader Header) *Message {
	m := createBaseMessage()
	m.Cmd = CmdReqHeaderValidation
	m.Data = blockHeader.serialize()
	return m
}

func createMsResponseHeaderValidation(isValid bool) *Message {
	m := createBaseMessage()
	m.Cmd = CmdResHeaderValidation
	m.Data = []byte(strconv.FormatBool(isValid))
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
