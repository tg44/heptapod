package cli_utils

import (
	"fmt"
	"github.com/tg44/heptapod/pkg"
	"github.com/tg44/heptapod/pkg/tmutil"
	"github.com/urfave/cli/v2"
	"log"
)

var file string
var all bool
var current bool

var TmCommands = &cli.Command{
	Name:    "timeMachine",
	Aliases: []string{"tm"},
	Usage:   "timeMachine related functions",
	Subcommands: []*cli.Command{
		TmListExcluded,
		TmListExcludedAll,
		TmPrune,
	},
}

var TmListExcluded = &cli.Command{
	Name:    "excluded",
	Aliases: []string{"ls", "list"},
	Usage:   "lists the excluded dirs from tmutil (not full, uses cache)",
	Flags:   []cli.Flag{},
	Action: func(c *cli.Context) error {
		res := tmutil.GetExcludeList()
		fmt.Println(res)
		return nil
	},
}

var TmListExcludedAll = &cli.Command{
	Name:    "excludedAll",
	Aliases: []string{},
	Usage:   "lists all the excluded dirs (not possible ATM)",
	Flags:   []cli.Flag{},
	Action: func(c *cli.Context) error {
		res := tmutil.GetExcludeListAll()
		fmt.Println(res)
		return nil
	},
}

var TmPrune = &cli.Command{
	Name:    "prune",
	Aliases: []string{},
	Usage:   "removes excludes from TM",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "all",
			Aliases:     []string{"a"},
			Value:       false,
			Usage:       "remove all previously added exclude paths from the log dir",
			Destination: &all,
		},
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Value:       "",
			Usage:       "remove previously added exclude paths from the file in the log dir",
			Destination: &file,
		},
		&cli.BoolFlag{
			Name:        "run",
			Aliases:     []string{"r"},
			Value:       false,
			Usage:       "remove previously added exclude paths based on the current rules",
			Destination: &current,
		},
	},
	Action: func(c *cli.Context) error {
		if current {
			res := pkg.GetExcludedPaths(rulePath, par, buffer, verbose)
			tmutil.RemovePathsFromTM(res, buffer, verbose)
		} else if all {
			tmutil.RemoveAllFromLogs(logDir, buffer, verbose)
		} else if file != "" {
			tmutil.RemoveFileFromLogs(logDir, file, buffer, verbose)
		} else {
			log.Fatal("one of the options is mandatory, please add current, all, or a file")
		}
		return nil
	},
}
