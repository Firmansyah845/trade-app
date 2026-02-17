package main

import (
	"os"

	"awesomeProjectCr/cmd"
	"awesomeProjectCr/cmd/app"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	Version   = "development"
	BuildTime = "unknown"
)

func main() {
	app.Init()
	defer app.Shutdown()

	cliApp := cli.NewApp()
	cliApp.Name = "awesomeProjectCr"
	cliApp.Version = Version
	cliApp.Metadata = map[string]interface{}{
		"buildTime": BuildTime,
	}

	cliApp.Commands = cli.Commands{
		{
			Name:  "server",
			Usage: "Start REST server",
			Action: func(c *cli.Context) error {
				cmd.StartServer()
				return nil
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run application")
	}
}
