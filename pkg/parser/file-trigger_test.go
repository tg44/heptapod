package parser

import "testing"

func TestFileTriggerSettingsParse(t *testing.T) {
	rule, err := ruleParse("../../test-rules/git.yaml")
	if err != nil {
		t.Errorf("Rule parser can't parse testRule!")
	}
	settings, err2 := fileTriggerSettingsParse(rule.RuleSettings)
	//t.Log(rule.RuleSettings)
	//t.Log(settings)
	if err2 != nil {
		t.Errorf("FileSettings parser can't parse testRule!")
	}
	if settings == nil {
		t.Errorf("FileSettings returns nil without error!")
	}
	if len(settings.ExcludePaths) != 1 || settings.ExcludePaths[0] != "." {
		t.Errorf("FileSettings parser ExcludePaths parser error!")
	}
	if settings.FileTrigger != ".git" {
		t.Errorf("FileSettings parser FileTrigger parser error!")
	}
}
