package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tg44/heptapod/pkg"
	"github.com/tg44/heptapod/pkg/parser"
	"github.com/tg44/heptapod/pkg/tmutil"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"strings"
)

var version string = "local-build"
var commit string = ""
var date string = ""

func main() {
	var rulePath string
	var logDir string
	var file string
	var all bool
	var dry bool
	var current bool
	var verbose bool
	var par int
	var buffer int

	app := &cli.App{
		Name:  "heptapod",
		Usage: "Fine-tune your TimeMachine excludes!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "rules",
				Aliases:     []string{"r"},
				Value:       "~/.heptapod/rules",
				Usage:       "the directory containing rule yamls",
				Destination: &rulePath,
			},
			&cli.StringFlag{
				Name:        "logDir",
				Aliases:     []string{"ld"},
				Value:       "~/.heptapod/logs",
				Usage:       "the directory where excluded dirs logged for reliable prune/uninstall",
				Destination: &logDir,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Value:       false,
				Usage:       "prints out performance logs (and other logs in general)",
				Destination: &verbose,
			},
			&cli.IntFlag{
				Name:        "parallelism",
				Aliases:     []string{"p", "par"},
				Value:       4,
				Usage:       "number of workers where the code is multithreaded",
				Destination: &par,
			},
			&cli.IntFlag{
				Name:        "bufferSize",
				Aliases:     []string{"b", "buffer"},
				Value:       2048,
				Usage:       "number of elements buffered, can cause deadlocks",
				Destination: &buffer,
			},
		},
		Commands: []*cli.Command{
			{
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
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list the rules",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Aliases:     []string{"a"},
						Value:       false,
						Usage:       "show all tables",
						Destination: &all,
					},
				},
				Action: func(c *cli.Context) error {
					return ruleTable(rulePath, all)
				},
			},
			{
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
						fmt.Print(strings.Join(res,"\r\n"))
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
			},
			{
				Name:    "excluded",
				Aliases: []string{"lse"},
				Usage:   "lists all the excluded dirs from tmutil",
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					res := tmutil.GetExcludeList()
					fmt.Println(res)
					return nil
				},
			},
			{
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
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func ruleTable(path string, all bool) error {
	list, err := parser.ParseRulesFromDir(path)
	if err != nil {
		return err
	}
	fmt.Println("Enabled rules:")
	writeRules(list.Enabled)
	if all {
		fmt.Println("Disabled rules:")
		writeRules(list.Disabled)
	}
	if all {
		fmt.Println("Errors:")
		writeErrorRules(list.FileErrors)
	}
	if all {
		fmt.Println("Type parse errors:")
		writeTypeErrorRules(list.TypeErrors)
	}
	return nil
}

func writeErrorRules(paths []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"nonparseable files"})
	table.SetBorder(false)

	for _, row := range paths {
		table.Append([]string{row})
	}

	table.SetAutoMergeCells(false)
	table.Render()
}

func writeTypeErrorRules(tes map[string]parser.Rule) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"path", "name", "type"})
	table.SetBorder(false)

	for k, v := range tes {
		table.Append([]string{k, v.Name, v.RuleType})
	}

	table.SetAutoMergeCells(false)
	table.Render()
}

func writeRules(tes map[string]parser.Rule) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"path", "name", "type", "search", "ignore"})
	table.SetBorder(false)

	for k, v := range tes {
		table.Append([]string{k, v.Name, v.RuleType, strings.Join(v.SearchPaths, ", "), strings.Join(v.IgnorePaths, ", ")})
	}

	table.SetAutoMergeCells(false)
	table.Render()
}
