package pkg

import (
	"github.com/tg44/heptapod/pkg/parser"
	"github.com/tg44/heptapod/pkg/utils"
	"github.com/tg44/heptapod/pkg/walker"
	"log"
)

func GetExcludedPaths(ruleDir string, par int, bufferSize int, verbose int) []string {
	path, err := utils.FixupPathsToHandleHome(ruleDir)
	if err != nil {
		log.Fatal(err)
	}
	jobs, err := parser.ParseFromDir(path)
	if err != nil {
		log.Fatal(err)
	}
	res := walker.Run(jobs, par, bufferSize, verbose)
	return res
}
