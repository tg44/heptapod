package parser

import (
	"testing"
)

func TestRuleParse(t *testing.T) {
	res, err := ruleParse("../test-rules/git.yaml")
	//t.Log(res)
	if err != nil {
		t.Errorf("Rule parser can't parse testRule!")
	}
	if res == nil {
		t.Errorf("Rule parser returned with a nil but not with an error!")
	}
	if res.Name != "git" {
		t.Errorf("Rule parser name parse error!")
	}
	if res.Enabled != false {
		t.Errorf("Rule parser enabled parse error!")
	}
	if len(res.SearchPaths) != 1 || res.SearchPaths[0] != "~" {
		t.Errorf("Rule parser SearchPaths parse error!")
	}
	if len(res.IgnorePaths) != 0 {
		t.Errorf("Rule parser IgnorePaths parse error!")
	}
	if res.RuleType != "file-trigger" {
		t.Errorf("Rule parser name parse error!")
	}
}
