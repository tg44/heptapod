package cli_utils

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func VersionCommand(version string, commit string, date string) *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "version info",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Value:       false,
				Usage:       "more detaild version info",
				Destination: &verbose,
			},
		},
		Action: func(c *cli.Context) error {
			if verbose {
				fmt.Println("version: ", version)
				fmt.Println("commit: ", commit)
				fmt.Println("date: ", date)
			} else {
				fmt.Println(version)
			}
			return nil
		},
	}
}
