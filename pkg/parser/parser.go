package parser

import (
	"github.com/tg44/heptapod/pkg/walker"
	"log"
)

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
	if rule.RuleType == "file-trigger" {
		settings, err2 := fileTriggerSettingsParse(rule.RuleSettings)
		if err2 != nil {
			log.Println(err2)
			return []walker.WalkJob{}
		}
		tasks := []walker.WalkJob{}
		walkerFun := fileTriggerWalker(*rule, *settings)
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
	for k, _ := range paths {
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
