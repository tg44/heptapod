package parser

import (
	"errors"
	"fmt"
	"github.com/tg44/heptapod/pkg/utils"
	"github.com/tg44/heptapod/pkg/walker"
	"os"
	"path/filepath"
)

type FileTriggerSettings struct {
	FileTrigger       string
	ExcludePaths []string
}

const FileTriggerType = ""

func fileTriggerSettingsParse(i map[string]interface{}) (*FileTriggerSettings, error) {
	fileTrigger, found1 := i["fileTrigger"].(string)
	excludePathsI, found2 := i["excludePaths"].([]interface{})
	excludePaths := make([]string, len(excludePathsI))
	for i, v := range excludePathsI {
		excludePaths[i] = fmt.Sprint(v)
	}
	if found1 && found2 {
		return &FileTriggerSettings{fileTrigger, excludePaths}, nil
	}
	return nil, errors.New("The given input can't be parsed to FileTriggerSettings!")
}

func fileTriggerWalker(rule Rule, settings FileTriggerSettings) walker.Walker {
	fixIgnorePaths := []string{}
	for _, p := range rule.IgnorePaths {
		path, err := utils.FixupPathsToHandleHome(p)
		if err == nil {
			fixIgnorePaths = append(fixIgnorePaths, path)
		}
	}
	return func(path string, subfiles []os.FileInfo) ([]string, []string, []string) {
		if utils.ContainsFIA(subfiles, settings.FileTrigger) {
			var ret []string
			for _, f := range settings.ExcludePaths {
				ret = append(ret, filepath.Join(path, f))
			}
			return ret, []string{}, fixIgnorePaths
		}
		return []string{}, []string{}, fixIgnorePaths
	}
}
