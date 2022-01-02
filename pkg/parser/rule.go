package parser

import (
	"github.com/tg44/heptapod/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Rule struct {
	Name       string `yaml:"name"`
	Enabled bool `yaml:"enabled"`
	SearchPaths []string `yaml:"searchPaths"`
	IgnorePaths []string `yaml:"ignorePaths"`
	RuleType string `yaml:"ruleType"`
	RuleSettings map[string]interface{} `yaml:"ruleSettings"`
}

func RuleParse(fileName string) (*Rule, error) {
	yfile, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	var rule Rule

	err2 := yaml.Unmarshal(yfile, &rule)

	if err2 != nil {
		return nil, err
	}

	return &rule, nil
}

func RuleWrite(rule Rule, file string) error {

	rule.IgnorePaths = utils.Unique(rule.IgnorePaths)
	data, err := yaml.Marshal(&rule)
	if err != nil {
		return err
	}

	err2 := ioutil.WriteFile(file, data, 0)
	if err2 != nil {
		return err2
	}
	return nil
}
