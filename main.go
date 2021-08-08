package main

import (
	"fmt"
	"os"

	"github.com/s3cmd-go/config"
	"github.com/urfave/cli/v2"
)

func CmdNotImplemented(*config.Config, *cli.Context) error {
	return fmt.Errorf("Command not implemented")
}

func main() {
	// Now the setup for the application

	cliapp := cli.NewApp()
	cliapp.Name = "s3cmd-go"
	// cliapp.Usage = ""
	cliapp.Version = "1.0.0"

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version, V",
		Usage: "print version number",
	}

	cliapp.Flags = BuildFlags()

	cliapp.Commands = BuildCommands(cliapp)

	cliapp.Run(os.Args)
}
