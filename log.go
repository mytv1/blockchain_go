package main

import (
	"io"
	"log"
)

var (
	// Trace log process's trace
	Trace *log.Logger

	// Info log program's state
	Info *log.Logger

	// Debug log debug mesage
	Debug *log.Logger

	// Warning log unwanted usage
	Warning *log.Logger

	// Error log errors
	Error *log.Logger
)

func initLog(
	traceHandle io.Writer,
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ltime|log.Lshortfile)

	Debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ltime|log.Lshortfile)
}
