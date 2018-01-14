package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func startServer(bc *Blockchain) {
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

		go handleRequest(conn, bc)
	}
}

func handleRequest(conn net.Conn, bc *Blockchain) {
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		Error.Println("Error reading:", err.Error())
		return
	}

	m := deserializeMessage(buf[:length])

	Info.Printf("Handle command %s request from : %s\n", m.Cmd, conn.RemoteAddr())

	switch m.Cmd {
	case CmdReqBestHeight:
		handleReqBestHeight(conn, bc)
	case CmdReqBlock:
		handleReqBlock(conn, bc, m)
	case CmdPrintBlockchain:
		handlePrintBlockchain(bc)
	case CmdReqAddBlock:
		handleReqAddBlock(conn, bc, m)
	case CmdSpreadHashList:
		handleSpreadHashList(conn, bc, m)
	default:
		Info.Printf("Message command invalid\n")
	}

	conn.Close()
}

func handleReqBestHeight(conn net.Conn, bc *Blockchain) {
	responseMs := createMsReponseBestHeight(bc.getBestHeight())
	conn.Write(responseMs.serialize())
}

func handleReqBlock(conn net.Conn, bc *Blockchain, m *Message) {
	number, err := strconv.Atoi(string(m.Data))
	if err != nil {
		Error.Printf("")
	}
	index := number - 1
	block := bc.Blocks[index]
	responseMs := createMsResponseBlock(block)
	conn.Write(responseMs.serialize())
}

func handlePrintBlockchain(bc *Blockchain) {
	Info.Printf("\n%v", bc)
}

func handleReqAddBlock(conn net.Conn, bc *Blockchain, m *Message) {
	bc.addBlock(string(m.Data))
	spreadHashList(bc)
}

func handleSpreadHashList(conn net.Conn, bc *Blockchain, m *Message) {
	Info.Printf("Blockchain's change detected. Start sync.")
	sendRequestBc(m.Source, bc)
}
