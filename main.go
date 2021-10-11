package main

import (
	"fmt"
	"github.com/tg44/heptapod/pkg/tmutil"
)

func main() {
	//res := pkg.GetExcludedPaths([]string{"test-rules/scala.yaml", "test-rules/node.yaml"}, 4, 2048)
	//tmutil.AddPathsToTM(res, 2048)
	//fmt.Println(len(res))

	fmt.Println(tmutil.GetExcludeList())
}
