package main

import (
	cli_utils "github.com/tg44/heptapod/pkg/cli-utils"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

var version string = "local-build"
var commit string = ""
var date string = ""

func main() {
	app := &cli.App{
		Name:  "heptapod",
		Usage: "Fine-tune your TimeMachine excludes!",
		Flags: []cli.Flag{
			cli_utils.RulesPathFlag,
			cli_utils.LogDirFlag,
			cli_utils.VerboseFlag,
			cli_utils.ParallelismFlag,
			cli_utils.BufferSizeFlag,
		},
		Commands: []*cli.Command{
			cli_utils.VersionCommand(version, commit, date),
			cli_utils.RuleCommands,
			cli_utils.RunCommand,
			cli_utils.TmCommands,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
