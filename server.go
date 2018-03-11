package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/google/go-cmp/cmp"
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

	m := new(Message)
	err = json.Unmarshal(buf[:length], m)

	if err != nil {
		Error.Println("Error unmarshal:", err.Error())
		return
	}

	Info.Printf("Handle command %s request from : %s\n", m.Cmd, conn.RemoteAddr())

	switch m.Cmd {
	case CmdReqBestHeight:
		handleReqBestHeight(conn, bc)
	case CmdReqBlock:
		handleReqBlock(conn, bc, m)
	case CmdPrintBlockchain:
		handlePrintBlockchain(bc)
	case CmdSpreadHashList:
		handleSpreadHashList(conn, bc, m)
	case CmdReqHeaderValidation:
		handleReqHeaderValidation(conn, bc, m)
	case CmdReqAddress:
		handleReqAddress(conn, m)
	case CmdReqPrintAllAddressInfo:
		handleReqPrintAllAddressInfo(conn, bc, m)
	case CmdReqAddTransaction:
		handleReqAddTransaction(conn, bc, m)
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
	requestBlockHeight := bytesToInt(m.Data)
	block := bc.getBlockByHeight(requestBlockHeight)
	responseMs := createMsResponseBlock(block)
	conn.Write(responseMs.serialize())
}

func handlePrintBlockchain(bc *Blockchain) {
	Info.Printf("\n%s", bc)
}

func handleReqAddTransaction(conn net.Conn, bc *Blockchain, m *Message) {
	var isSuccess bool
	tx := deserializeTransaction(m.Data)
	Info.Printf("Receive tx : %s", tx)
	isSuccess = bc.verifyTransaction(tx)
	if isSuccess == true {
		Info.Printf("Transaction is valid. Create new block.")
		Info.Printf("Append coinbase transaction to address %s.", getWallet().Address)
		coinbaseTx := newCoinbaseTx(getWallet().Address)
		newBlock := newBlock([]Transaction{*tx, *coinbaseTx}, bc.getTopBlockHash(), bc.getBestHeight()+1)
		bc.addBlock(newBlock)
		spreadHashList(bc)
	} else {
		Info.Printf("Transaction is invalid. Nothing happened.")
	}
	responseMs := createMsResponseAddTransaction(isSuccess)
	conn.Write(responseMs.serialize())
}

func handleSpreadHashList(conn net.Conn, bc *Blockchain, m *Message) {
	Info.Printf("Blockchain's change detected. Start sync.")
	sendRequestBc(m.Source, bc)
}

func handleReqHeaderValidation(conn net.Conn, bc *Blockchain, m *Message) {
	oppHeader := deserializeHeader(m.Data)
	myBlock := bc.getBlockByHeight(oppHeader.Height)
	result := cmp.Equal(*oppHeader, myBlock.Header)
	responseMs := createMsResponseHeaderValidation(result)
	conn.Write(responseMs.serialize())
}

func handleReqAddress(conn net.Conn, m *Message) {
	responseMs := createMsResponseAddress()
	conn.Write(responseMs.serialize())
	Info.Printf("My Address : %s", getConfig().SWallet.Address)
}

func handleReqPrintAllAddressInfo(conn net.Conn, bc *Blockchain, m *Message) {
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()
	addressInfos := UTXOSet.getAllAddressInfo()

	Info.Print(" All address information 	")
	Info.Printf("|             Address                | Value |")
	for pubKeyHashAsStr, val := range addressInfos {
		pubKeyHash, err := hex.DecodeString(pubKeyHashAsStr)
		if err != nil {
			Error.Fatal(err.Error())
		}
		Info.Printf("| %s | %5d |", generateAddress(pubKeyHash), val)
	}

}
