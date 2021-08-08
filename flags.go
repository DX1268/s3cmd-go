package main

import "github.com/urfave/cli/v2"

func BuildFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "config, c",
			Value: cli.NewStringSlice("$HOME/.s3cfg"),
			Usage: "Config `FILE` name.",
		},
		&cli.StringFlag{
			Name:    "access-key",
			Usage:   "AWS Access Key `ACCESS_KEY`",
			EnvVars: []string{"AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY"},
		},
		&cli.StringFlag{
			Name:    "secret-key",
			Usage:   "AWS Secret Key `SECRET_KEY`",
			EnvVars: []string{"AWS_SECRET_ACCESS_KEY", "AWS_SECRET_KEY"},
		},
		&cli.StringFlag{
			Name:    "storage-class",
			Usage:   "Storage class (default: STANDARD)",
			EnvVars: []string{"AWS_S3_STORAGE_CLASS"},
		},
		&cli.IntFlag{
			Name:    "concurrency",
			Usage:   "Concurrency `NUM`",
			EnvVars: []string{"AWS_S3_CONCURRENCY"},
			Value:   5,
		},
		&cli.IntFlag{
			Name:    "part-size",
			Usage:   "Part size `NUM` in mb",
			EnvVars: []string{"AWS_S3_PARTSIZE"},
			Value:   5,
		},
		&cli.BoolFlag{
			Name:  "recursive,r",
			Usage: "Recursive upload, download or removal",
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force overwrite and other dangerous operations.",
		},
		&cli.BoolFlag{
			Name:  "skip-existing",
			Usage: "Skip over files that exist at the destination (only for [get] and [sync] commands).",
		},
		&cli.BoolFlag{
			Name:  "verbose,v",
			Usage: "Verbose output (e.g. debugging)",
		},
		&cli.BoolFlag{
			Name:  "dry-run,n",
			Usage: "Only show what should be uploaded or downloaded but don't actually do it. May still perform S3 requests to get bucket listings and other information though (only for file transfer commands)",
		},
		&cli.BoolFlag{
			Name:  "check-md5",
			Usage: "Check MD5 sums when comparing files for [sync]. (not default)",
		},
		&cli.BoolFlag{
			Name:  "no-check-md5",
			Usage: "Do not check MD5 sums when comparing files for [sync] (default).",
		},
	}
}
