package parser

import (
	"errors"
	"github.com/tg44/heptapod/pkg/utils"
	"github.com/tg44/heptapod/pkg/walker"
	"os"
	"path/filepath"
)

type GlobalSettings struct {
	Path       string
	IsExcluded bool
	IsIgnored  bool
	/**
	excluded - ignored - result
	t - t/f - add to exclude and globally ignore from rules
	f - t - globally ignore from rules
	f - f - do nothing
	*/

}

func globalSettingsParse(i map[string]interface{}) (*GlobalSettings, error) {
	path, found1 := i["path"].(string)
	handleWith, found2 := i["handleWith"].(string)
	var isEx = false
	var isIgn = false
	if handleWith == "exclude" {
		isEx = true
	} else if handleWith == "ignore" {
		isIgn = true
	}
	if found1 && found2 {
		return &GlobalSettings{path, isEx, isIgn}, nil
	}
	return nil, errors.New("The given input can't be parsed to GlobalSettings!")
}

func getGlobalIgnore(settings GlobalSettings) []string {
	fixIgnorePaths := []string{}
	if settings.IsIgnored || settings.IsExcluded {
		path, err := utils.FixupPathsToHandleHome(settings.Path)
		if err == nil {
			fixIgnorePaths = append(fixIgnorePaths, path)
		}
	}
	return fixIgnorePaths
}

func globalWalker(rule Rule, settings GlobalSettings) walker.Walker {
	fixIgnorePaths := getGlobalIgnore(settings)
	fixExcludePaths := []string{}
	if settings.IsExcluded {
		fixExcludePaths = fixIgnorePaths
	}

	return func(path string, subfiles []os.FileInfo) ([]string, []string, []string, string) {
		var localIgnore []string
		for _, f := range subfiles {
			localIgnore = append(localIgnore, filepath.Join(path, f.Name()))
		}

		return fixExcludePaths, fixIgnorePaths, localIgnore, rule.Name
	}
}
