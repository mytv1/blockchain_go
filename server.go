package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func StartServer() {
	config = GetConfig()
	l, err := net.Listen("tcp", config.Nw.LocalNode.Address)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer l.Close()

	Info.Println("Node listening on " + config.Nw.LocalNode.Address)

	for {
		conn, err := l.Accept()
		if err != nil {
			Error.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		Error.Println("Error reading:", err.Error())
		return
	}

	bc := GetBlockchain()

	var m *Message = new(Message)
	err = json.Unmarshal(buf[:length], m)

	if err != nil {
		Error.Println("Error unmarshal:", err.Error())
		return
	}

	Info.Printf("Handle command %s request from : %s\n", m.Cmd, conn.RemoteAddr())

	switch m.Cmd {
	case CMD_REQ_BLOCKCHAIN:
		conn.Write(bc.SerializeBlockchain())
	case CMD_REQ_BEST_HEIGHT:
		responseMs := CreateMsReponseBestHeight(bc.GetBestHeight())
		conn.Write(responseMs.Serialize())
	case CMD_REQ_BLOCK:
		block := bc.Blocks[uint8(m.Data[0])-1]
		responseMs := CreateMsResponseBlock(block)
		conn.Write(responseMs.Serialize())
	case CMD_PRINT_BLOCKCHAIN:
		Info.Printf("\n%v", bc)
	case CMD_REQ_ADD_BLOCK:
		bc.AddBlock(string(m.Data))
		SpreadHashList()
	case CMD_SPREAD_HASHLIST:
		Info.Printf("Blockchain's change detected. Start sync.")
		SendRequestBc(m.Source, bc)
	default:
		Info.Printf("Message command invalid\n")
	}

	conn.Close()
}
