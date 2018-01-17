package main

import (
	"io/ioutil"
	"os"
)

func main() {
	initLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	app := newCliApp()
	app.Run(os.Args)
}
