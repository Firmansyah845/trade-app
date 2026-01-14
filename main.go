package main

import (
	"os"

	"awesomeProjectCr/cmd"
	"awesomeProjectCr/cmd/app"

	"github.com/urfave/cli/v2"
)

func main() {
	app.Init()
	defer app.Shutdown()

	cliApp := cli.NewApp()
	cliApp.Name = "awesomeProjectCr"

	cliApp.Commands = cli.Commands{
		{
			Name:  "server",
			Usage: "Start server",
			Action: func(c *cli.Context) error {
				cmd.StartServer()
				return nil
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
