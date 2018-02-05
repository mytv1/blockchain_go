package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "simple blockchain"
	app.Usage = "simple blockchain implemented by golang"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{}

	initCreateWalletCLI(app)
	initStartServerCLI(app)

	return app
}

func initCreateWalletCLI(app *cli.App) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    "createwallet",
		Aliases: []string{"cw"},
		Usage:   "start server",
		Action: func(c *cli.Context) error {
			config := initConfig(defaultConfigPath)
			wallet := newWallet()

			config.SWallet = *wallet.toStorable()
			config.exportConfig(defaultConfigPath)
			fmt.Printf("New wallet is created successfully! Wallet is exported to : * %s *\n", defaultConfigPath)
			fmt.Printf("  + Private Key : %s\n", config.SWallet.PrivateKey)
			fmt.Printf("  + Public Key : %s\n", config.SWallet.PublicKey)
			fmt.Printf("  + Address : %s\n", config.SWallet.Address)
			return nil
		},
	})
}

func initStartServerCLI(app *cli.App) {
	var configPath string
	app.Flags = append(app.Flags, cli.StringFlag{
		Name:        "config, c",
		Value:       defaultConfigPath,
		Usage:       "Load configuration form `FILE`",
		Destination: &configPath,
	})

	app.Commands = append(app.Commands, cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start server",
		Action: func(c *cli.Context) error {
			execStartCmd(c, configPath)
			return nil
		},
	})
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
