package main

import (
	"github.com/urfave/cli"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "simple blockchain"
	app.Usage = "simple blockchain implemented by golang"

	initStartServerCLI(app)

	return app
}

func initStartServerCLI(app *cli.App) {
	var configPath string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       defaultConfigPath,
			Usage:       "Load configuration form `FILE`",
			Destination: &configPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start server",
			Action: func(c *cli.Context) error {
				execStartCmd(c, configPath)
				return nil
			},
		},
	}
}

func execStartCmd(c *cli.Context, configPath string) {
	initConfig(configPath)
	bc := getLocalBc()
	if bc == nil {
		Info.Printf("Local blockchain database not found. Create new empty blockchain (size = 0).")
		bc = createEmptyBlockchain()
	} else {
		Info.Printf("Read blockchain from local database completed.")
	}
	syncWithNeighborNode(bc)

	if bc.isEmpty() {
		Info.Printf("No avaiable node for synchronization. Init new blockchain.")
		bc.addBlock(newGenesisBlock())
	}

	startServer(bc)
	defer bc.db.Close()
}
