package cli_utils

import (
	"fmt"
	"github.com/tg44/heptapod/pkg"
	"github.com/tg44/heptapod/pkg/tmutil"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
)

var dry bool

var RunCommand = &cli.Command{
	Name:    "run",
	Aliases: []string{},
	Usage:   "run the exclude detection, and also exclude the dirs",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "dry",
			Aliases:     []string{"d"},
			Value:       false,
			Usage:       "only prints the paths we should exclude (if exists)",
			Destination: &dry,
		},
	},
	Action: func(c *cli.Context) error {
		if dry {
			res := pkg.GetExcludedPaths(rulePath, par, buffer, verbose)
			fmt.Println("-----")
			fmt.Print(strings.Join(res, "\r\n"))
			fmt.Print("\n")
		} else {
			res := pkg.GetExcludedPaths(rulePath, par, buffer, verbose)
			tmutil.AddPathsToTM(res, logDir, buffer, verbose)
			if verbose {
				log.Printf("total %d paths found\n", len(res))
			}
		}
		return nil
	},
}
