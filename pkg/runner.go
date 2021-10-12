package pkg

import (
	"github.com/tg44/heptapod/pkg/parser"
	"github.com/tg44/heptapod/pkg/walker"
	"log"
)

func GetExcludedPaths(ruleDir string, par int, bufferSize int, verbose bool) []string {
	jobs, err := parser.ParseFromDir(ruleDir)
	if err != nil {
		log.Fatal(err)
	}
	res := walker.Run(jobs, par, bufferSize, verbose)
	return res
}
