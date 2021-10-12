package parser

import (
	"github.com/tg44/heptapod/pkg/walker"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func ParseFromDir(path string) ([]walker.WalkJob, error) {
	files, err := ioutil.ReadDir(path)
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
	TypeErrors map[string]Rule
	Enabled    map[string]Rule
	Disabled   map[string]Rule
}

func ParseRulesFromDir(path string) (*RuleGroups, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	errors := []string{}
	typeErrors := map[string]Rule{}
	enabled := map[string]Rule{}
	disabled := map[string]Rule{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml") {
			fp := filepath.Join(path, file.Name())
			rule, err := ruleParse(fp)
			if err != nil {
				errors = append(errors, file.Name())
				continue
			}
			rt := parseRuleTypes(*rule)
			if len(rt) == 0 {
				typeErrors[file.Name()] = *rule
			} else if !rule.Enabled {
				disabled[file.Name()] = *rule
			} else {
				enabled[file.Name()] = *rule
			}
		}
	}

	return &RuleGroups{errors, typeErrors, enabled, disabled}, nil
}

func Parse(ruleFiles []string) []walker.WalkJob {
	jobs := []walker.WalkJob{}
	for _, f := range ruleFiles {
		jobs = append(jobs, parse(f)...)
	}
	jobs = mergeJobs(jobs)
	return jobs
}

func parse(ruleFile string) []walker.WalkJob {
	rule, err := ruleParse(ruleFile)
	if err != nil {
		log.Println(err)
		return []walker.WalkJob{}
	}
	if !rule.Enabled {
		return []walker.WalkJob{}
	}

	return parseRuleTypes(*rule)
}

func parseRuleTypes(rule Rule) []walker.WalkJob {
	if rule.RuleType == "file-trigger" {
		settings, err2 := fileTriggerSettingsParse(rule.RuleSettings)
		if err2 != nil {
			log.Println(err2)
			return []walker.WalkJob{}
		}
		tasks := []walker.WalkJob{}
		walkerFun := fileTriggerWalker(rule, *settings)
		for _, p := range rule.SearchPaths {
			tasks = append(tasks, walker.WalkJob{p, []walker.Walker{walkerFun}, []string{}})
		}
		return tasks
	}
	return []walker.WalkJob{}
}

func mergeJobs(works []walker.WalkJob) []walker.WalkJob {
	paths := map[string]bool{}
	for _, w := range works {
		paths[w.Rootpath] = true
	}
	pathArr := []string{}
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
