package parser

import (
	"errors"
	"fmt"
	"github.com/tg44/heptapod/pkg/utils"
	"github.com/tg44/heptapod/pkg/walker"
	"os"
)

type IgnoreTriggerSettings struct {
	FileName          string
	ForceIncludePaths []string
}

func ignoreTriggerSettingsParse(i map[string]interface{}) (*IgnoreTriggerSettings, error) {
	fileName, found1 := i["fileName"].(string)
	forceIncludePathsI, found2 := i["forceIncludePaths"].([]interface{})
	forceIncludePaths := make([]string, len(forceIncludePathsI))
	for i, v := range forceIncludePathsI {
		forceIncludePaths[i] = fmt.Sprint(v)
	}
	if found1 && found2 {
		return &IgnoreTriggerSettings{fileName, forceIncludePaths}, nil
	}
	return nil, errors.New("The given input can't be parsed to IgnoreTriggerSettings!")
}

func ignoreTriggerWalker(rule Rule, settings IgnoreTriggerSettings) walker.Walker {
	fixIgnorePaths := []string{}
	for _, p := range rule.IgnorePaths {
		path, err := utils.FixupPathsToHandleHome(p)
		if err == nil {
			fixIgnorePaths = append(fixIgnorePaths, path)
		}
	}

	return func(path string, subfiles []os.DirEntry) ([]string, []string, []string, string) {
		/*
			//here we should go in and handle all the subpaths by this given file
			//not as trivial as it seemed at first bcs we need to implement something similar as the walker

			if utils.ContainsFIA(subfiles, settings.FileName) {
				object, err := ignore.CompileIgnoreFile(filepath.Join(path, settings.FileName))
				var ret []string

				return ret, []string{}, fixIgnorePaths, rule.Name
			}
		*/
		return []string{}, []string{}, fixIgnorePaths, rule.Name
	}
}
