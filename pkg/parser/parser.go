package parser

import (
	"github.com/tg44/heptapod/pkg/utils"
	"github.com/tg44/heptapod/pkg/walker"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ParseFromDir(path string) ([]walker.WalkJob, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for _, file := range files {
		paths = append(paths, filepath.Join(path, file.Name()))
	}
	return Parse(paths), err
}

type RuleGroups struct {
	FileErrors []string
	TypeErrors []Rule
	Enabled    []Rule
	Disabled   []Rule
}

func ParseRulesFromDir(pathIn string) (*RuleGroups, error) {
	path, err := utils.FixupPathsToHandleHome(pathIn)
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	errors := []string{}
	typeErrors := []Rule{}
	enabled := []Rule{}
	disabled := []Rule{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml") {
			fp := filepath.Join(path, file.Name())
			rule, err := RuleParse(fp)
			if err != nil {
				errors = append(errors, file.Name())
				continue
			}
			rt, _, _ := parseRuleTypes(rule)
			if len(rt) == 0 {
				typeErrors = append(typeErrors, *rule)
			} else {
				for _, r := range rule.flattened() {
					if !r.Enabled {
						disabled = append(disabled, r)
					} else {
						enabled = append(enabled, r)
					}
				}
			}
		}
	}

	return &RuleGroups{errors, typeErrors, enabled, disabled}, nil
}

func Parse(ruleFiles []string) []walker.WalkJob {
	jobs := []walker.WalkJob{}
	globalIgnores := []string{}
	for _, f := range ruleFiles {
		jobArr, ignores := parse(f)
		jobs = append(jobs, jobArr...)
		globalIgnores = append(globalIgnores, ignores...)
	}
	jobs = mergeJobs(jobs, globalIgnores)
	return jobs
}

func parse(ruleFile string) ([]walker.WalkJob, []string) {
	rule, err := RuleParse(ruleFile)
	if err != nil {
		log.Println(err)
		return []walker.WalkJob{}, []string{}
	}
	if rule == nil {
		log.Println(ruleFile, " has a strange error, we didn't get parse error, but it is not parsed either... Please report the content of the file to https://github.com/tg44/heptapod/issues/5")
		return []walker.WalkJob{}, []string{}
	}
	if !rule.Enabled {
		return []walker.WalkJob{}, []string{}
	}

	w, i, _ := parseRuleTypes(rule)
	return w, i
}

func parseRuleTypes(rule *Rule) ([]walker.WalkJob, []string, []Rule) {
	if rule.RuleType == "file-trigger" {
		settings, err2 := fileTriggerSettingsParse(rule.RuleSettings)
		if err2 != nil {
			log.Println(err2)
			return []walker.WalkJob{}, []string{}, []Rule{}
		}
		tasks := []walker.WalkJob{}
		walkerFun := fileTriggerWalker(*rule, *settings)
		for _, p := range rule.SearchPaths {
			tasks = append(tasks, walker.WalkJob{p, []walker.Walker{walkerFun}, []string{}})
		}
		return tasks, []string{}, []Rule{}
	} else if rule.RuleType == "global" {
		settings, err2 := globalSettingsParse(rule.RuleSettings)
		if err2 != nil {
			log.Println(err2)
			return []walker.WalkJob{}, []string{}, []Rule{}
		}
		tasks := []walker.WalkJob{}
		walkerFun := globalWalker(*rule, *settings)
		tasks = append(tasks, walker.WalkJob{"/", []walker.Walker{walkerFun}, []string{}})
		return tasks, getGlobalIgnore(*settings), []Rule{}
	} else if rule.RuleType == "list" {
		settings, err2 := listSettingsParse(rule.RuleSettings, parseRuleTypes)
		if err2 != nil {
			log.Println(err2)
		}
		rule.SubRules = settings.SubRules
		return settings.walkers, settings.globalIgnores, settings.SubRules
	}
	return []walker.WalkJob{}, []string{}, []Rule{}
}

func mergeJobs(works []walker.WalkJob, globalIgnores []string) []walker.WalkJob {
	paths := map[string]bool{}
	for _, w := range works {
		paths[w.Rootpath] = true
	}
	pathArr := globalIgnores
	for k := range paths {
		pathArr = append(pathArr, k)
	}
	newJobs := []walker.WalkJob{}
	for _, p := range pathArr {
		walkers := []walker.Walker{}
		for _, w := range works {
			if w.Rootpath == p {
				walkers = append(walkers, w.Walkers...)
			}
		}
		newJobs = append(newJobs, walker.WalkJob{p, walkers, pathArr})
	}
	return newJobs
}
