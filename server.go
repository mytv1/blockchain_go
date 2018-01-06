package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func startServer() {
	config = getConfig()
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

	bc := getBlockchain()

	m := new(Message)
	err = json.Unmarshal(buf[:length], m)

	if err != nil {
		Error.Println("Error unmarshal:", err.Error())
		return
	}

	Info.Printf("Handle command %s request from : %s\n", m.Cmd, conn.RemoteAddr())

	switch m.Cmd {
	case CmdReqBlockchain:
		conn.Write(bc.serialize())
	case CmdReqBestHeight:
		responseMs := createMsReponseBestHeight(bc.getBestHeight())
		conn.Write(responseMs.serialize())
	case CmdReqBlock:
		block := bc.Blocks[uint8(m.Data[0])-1]
		responseMs := createMsResponseBlock(block)
		conn.Write(responseMs.serialize())
	case CmdPrintBlockchain:
		Info.Printf("\n%v", bc)
	case CmdReqAddBlock:
		bc.addBlock(string(m.Data))
		spreadHashList()
	case CmdSpreadHashList:
		Info.Printf("Blockchain's change detected. Start sync.")
		sendRequestBc(m.Source, bc)
	default:
		Info.Printf("Message command invalid\n")
	}

	conn.Close()
}
