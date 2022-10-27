package cli_utils

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tg44/heptapod/pkg/parser"
	"github.com/tg44/heptapod/pkg/utils"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
)

var RuleCommands = &cli.Command{
	Name:    "rules",
	Aliases: []string{},
	Usage:   "rule related functions",
	Subcommands: []*cli.Command{
		RuleList,
		RuleAddGlobalIgnore,
		RuleRemoveGlobalIgnore,
		RuleEnable,
		RuleDisable,
	},
}

var RuleAddGlobalIgnore = &cli.Command{
	Name:    "ignoreAdd",
	Aliases: []string{},
	Usage:   "add dir as ignored for all rules",
	Action: func(c *cli.Context) error {
		return ruleIgnoreAddAll(rulePath, c.Args().Slice())
	},
}

var RuleRemoveGlobalIgnore = &cli.Command{
	Name:    "ignoreRemove",
	Aliases: []string{},
	Usage:   "remove ignored dir from all rules",
	Action: func(c *cli.Context) error {
		return ruleIgnoreRemoveAll(rulePath, c.Args().Slice())
	},
}

var RuleEnable = &cli.Command{
	Name:    "enable",
	Aliases: []string{},
	Usage:   "enables the given rules",
	Action: func(c *cli.Context) error {
		return ruleEnable(rulePath, c.Args().Slice())
	},
}

var RuleDisable = &cli.Command{
	Name:    "disable",
	Aliases: []string{},
	Usage:   "disables the given rules",
	Action: func(c *cli.Context) error {
		return ruleDisable(rulePath, c.Args().Slice())
	},
}

var RuleList = &cli.Command{
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
}

func ruleEnable(pathIn string, enables []string) error {
	path, err := utils.FixupPathsToHandleHome(pathIn)
	if err != nil {
		fmt.Println("Path error.")
		return err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Dir read error.")
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml") {
			fp := filepath.Join(path, file.Name())
			rule, err := parser.RuleParse(fp)
			if err != nil {
				continue
			}
			if !rule.Enabled && utils.ContainsSA(enables, rule.Name) {
				rule.Enabled = true
				err2 := parser.RuleWrite(*rule, fp)
				if err2 != nil {
					fmt.Println("Rule write error: ", fp)
					continue
				}
				fmt.Println("Enable: ", fp)
			}
		}
	}
	return nil
}

func ruleDisable(pathIn string, enables []string) error {
	path, err := utils.FixupPathsToHandleHome(pathIn)
	if err != nil {
		fmt.Println("Path error.")
		return err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Dir read error.")
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml") {
			fp := filepath.Join(path, file.Name())
			rule, err := parser.RuleParse(fp)
			if err != nil {
				continue
			}
			if rule.Enabled && utils.ContainsSA(enables, rule.Name) {
				rule.Enabled = false
				err2 := parser.RuleWrite(*rule, fp)
				if err2 != nil {
					fmt.Println("Rule write error: ", fp)
					continue
				}
				fmt.Println("Disable: ", fp)
			}
		}
	}
	return nil
}

func ruleIgnoreAddAll(pathIn string, excludePaths []string) error {
	path, err := utils.FixupPathsToHandleHome(pathIn)
	if err != nil {
		fmt.Println("Path error.")
		return err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Dir read error.")
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml") {
			fp := filepath.Join(path, file.Name())
			rule, err := parser.RuleParse(fp)
			if err != nil {
				continue
			}
			rule.IgnorePaths = append(rule.IgnorePaths, excludePaths...)
			err2 := parser.RuleWrite(*rule, fp)
			if err2 != nil {
				fmt.Println("Rule write error: ", fp)
				continue
			}
		}
	}
	return nil
}

func ruleIgnoreRemoveAll(pathIn string, excludePaths []string) error {
	path, err := utils.FixupPathsToHandleHome(pathIn)
	if err != nil {
		fmt.Println("Path error.")
		return err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Dir read error.")
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml") {
			fp := filepath.Join(path, file.Name())
			rule, err := parser.RuleParse(fp)
			if err != nil {
				continue
			}
			for _, n := range excludePaths {
				rule.IgnorePaths = utils.Filter(rule.IgnorePaths, func(s string) bool { return s == n })
			}
			err2 := parser.RuleWrite(*rule, fp)
			if err2 != nil {
				fmt.Println("Rule write error: ", fp)
				continue
			}
		}
	}
	return nil
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

func writeTypeErrorRules(tes []parser.Rule) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"path", "name", "type"})
	table.SetBorder(false)

	for _, v := range tes {
		table.Append([]string{v.FileName, v.Name, v.RuleType})
	}

	table.SetAutoMergeCells(false)
	table.Render()
}

func writeRules(tes []parser.Rule) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"path", "name", "type", "search", "ignore"})
	table.SetBorder(false)

	for _, v := range tes {
		table.Append([]string{v.FileName, v.Name, v.RuleType, strings.Join(v.SearchPaths, ", "), strings.Join(v.IgnorePaths, ", ")})
	}

	table.SetAutoMergeCells(false)
	table.Render()
}
