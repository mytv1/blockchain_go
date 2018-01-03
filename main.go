package main

import (
	"io/ioutil"
	"os"
)

func main() {
	InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	app := NewCliApp()
	app.Run(os.Args)
}
