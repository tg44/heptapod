package cli_utils

import (
	"github.com/urfave/cli/v2"
)

var rulePath string
var logDir string
var verbose int
var par int
var buffer int

var RulesPathFlag = &cli.StringFlag{
	Name:        "rules",
	Aliases:     []string{"r"},
	Value:       "~/.heptapod/rules",
	Usage:       "the directory containing rule yamls",
	Destination: &rulePath,
}

var LogDirFlag = &cli.StringFlag{
	Name:        "logDir",
	Aliases:     []string{"ld"},
	Value:       "~/.heptapod/logs",
	Usage:       "the directory where excluded dirs logged for reliable prune/uninstall",
	Destination: &logDir,
}

var VerboseFlag = &cli.IntFlag{
	Name:        "verbose",
	Aliases:     []string{"v"},
	Value:       0,
	Usage:       "prints out performance logs (and other logs in general) (0-4)",
	Destination: &verbose,
}

var ParallelismFlag = &cli.IntFlag{
	Name:        "parallelism",
	Aliases:     []string{"p", "par"},
	Value:       4,
	Usage:       "number of workers where the code is multithreaded",
	Destination: &par,
}

var BufferSizeFlag = &cli.IntFlag{
	Name:        "bufferSize",
	Aliases:     []string{"b", "buffer"},
	Value:       2048,
	Usage:       "number of elements buffered, can cause deadlocks",
	Destination: &buffer,
}
