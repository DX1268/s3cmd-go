package main

import (
	"fmt"

	"github.com/s3cmd-go/config"
	"github.com/s3cmd-go/controllers"
	"github.com/urfave/cli/v2"
)

type CmdHandler func(*config.Config, *cli.Context) error

// The wrapper to launch a command -- take care of standard setup
//  before we get going
func BuildLaunch(handler CmdHandler) func(*cli.Context) error {
	return func(c *cli.Context) error {
		config, err := config.NewConfig(c)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if err := handler(config, c); err != nil {
			fmt.Println(err)
			return err
		}
		return err
	}
}

func BuildCommands(cliapp *cli.App) []*cli.Command {
	return []*cli.Command{
		{
			Name:   "mb",
			Usage:  "Make bucket -- s3cmd-go mb s3://BUCKET",
			Action: BuildLaunch(controllers.MakeBucket),
			Flags:  cliapp.Flags,
		},
		{
			Name:   "rb",
			Usage:  "Remove bucket -- s3cmd-go mb s3://BUCKET",
			Action: BuildLaunch(controllers.RemoveBucket),
			Flags:  cliapp.Flags,
		},
		{
			Name:   "ls",
			Usage:  "List objects or buckets -- s3cmd-go ls [s3://BUCKET[/PREFIX]]",
			Action: BuildLaunch(controllers.ListBucket),
			Flags:  cliapp.Flags,
		},
		{
			Name:   "la",
			Usage:  "List all object in all buckets -- s3cmd-go la",
			Action: BuildLaunch(controllers.ListAll),
			Flags:  cliapp.Flags,
		},
		{
			Name:  "put",
			Usage: "Put file from bucket (really 'cp') -- s3cmd-go put FILE [FILE....] s3://BUCKET/PREFIX",
			//Action: BuildLaunch(CmdCopy),
			Flags: cliapp.Flags,
		},
		{
			Name:  "get",
			Usage: "Get file from bucket (really 'cp') -- s3cmd-go get s3://BUCKET/OBJECT LOCAL_FILE",
			//Action: BuildLaunch(CmdCopy),
			Flags: cliapp.Flags,
		},
		{
			Name:  "del",
			Usage: "Delete file from bucket -- s3cmd-go del s3://BUCKET/OBJECT",
			//Action: BuildLaunch(DeleteObjects),
			Flags: cliapp.Flags,
		},
		{
			Name:  "rm",
			Usage: "Delete file from bucket (del synonym) -- s3cmd-go rm s3://BUCKET/OBJECT",
			//Action: BuildLaunch(DeleteObjects),
			Flags: cliapp.Flags,
		},
		{
			Name:  "du",
			Usage: "Disk usage by buckets -- [s3://BUCKET[/PREFIX]]",
			//Action: BuildLaunch(GetUsage),
			Flags: cliapp.Flags,
		},
		{
			Name:  "cp",
			Usage: "copy files and directories -- SRC [SRC...] DST",
			//Action: BuildLaunch(CmdCopy),
			Flags: cliapp.Flags,
		},
		{
			Name:  "sync",
			Usage: "Synchronize a directory tree to S3 -- LOCAL_DIR s3://BUCKET[/PREFIX] or s3://BUCKET[/PREFIX] LOCAL_DIR",
			//Action: BuildLaunch(CmdSync),
			Flags: cliapp.Flags,
		},
		{
			Name:  "modify",
			Usage: "Modify object metadata -- s3://BUCKET1/OBJECT",
			//Action: BuildLaunch(Modify),
			Flags: cliapp.Flags,
		},
		{
			Name:  "info",
			Usage: "Get various information about Buckets or Files -- s3://BUCKET[/OBJECT]",
			//Action: BuildLaunch(GetInfo),
			Flags: cliapp.Flags,
		},
		// info
		// mv
	}
}
