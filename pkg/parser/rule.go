package parser

import (
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

func ruleParse(fileName string) (*Rule, error) {
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
