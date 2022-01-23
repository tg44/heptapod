package parser

import (
	"errors"
	"github.com/tg44/heptapod/pkg/walker"
	"gopkg.in/yaml.v2"
	"log"
)

type ListSettings struct {
	SubRules      []Rule
	walkers       []walker.WalkJob
	globalIgnores []string
}

func listSettingsParse(i map[string]interface{}, parser func(*Rule) ([]walker.WalkJob, []string, []Rule)) (*ListSettings, error) {
	subsI, found1 := i["subRules"].([]interface{})
	walkers := []walker.WalkJob{}
	ignores := []string{}
	rules := []Rule{}
	for _, v := range subsI {
		var rule Rule
		data, err := yaml.Marshal(&v)
		if err != nil {
			log.Println("subrule cannot be marshalled back")
			continue
		}
		err2 := yaml.Unmarshal(data, &rule)
		if err2 != nil {
			log.Println("subrule cannot be parsed back as a rule")
		}
		w, i, _ := parser(&rule)
		walkers = append(walkers, w...)
		ignores = append(ignores, i...)
		rules = append(rules, rule)
	}
	if found1 {
		return &ListSettings{rules, walkers, ignores}, nil
	}
	return nil, errors.New("The given input can't be parsed to ListSettings!")
}
