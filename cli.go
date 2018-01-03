package main

import (
	"github.com/urfave/cli"
)

func NewCliApp() *cli.App {
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
			Value:       DEFAULT_CONFIG_PATH,
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
	InitConfig(configPath)
	GetNeighborBc()
	StartServer()
}
