package pkg

import (
	"github.com/tg44/heptapod/pkg/parser"
	"github.com/tg44/heptapod/pkg/walker"
)

func GetExcludedPaths(ruleFiles []string, par int, bufferSize int) []string {
	jobs := parser.Parse(ruleFiles)
	res := walker.Run(jobs, par, bufferSize)
	return res
}
