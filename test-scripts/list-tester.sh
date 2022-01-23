mkdir -p ~/test_hepta/rules
cat << EOF > ~/test_hepta/rules/test-rule.yaml
name: "node-multi"
enabled: true
searchPaths: []
ignorePaths: []
ruleType: "list"
ruleSettings:
  subRules:
    - name: "test-node"
      enabled: true
      searchPaths: []
      ignorePaths: []
      ruleType: "global"
      ruleSettings:
        path: "~/.npm"
        handleWith: "exclude"
    - name: "test-node"
      enabled: true
      searchPaths: ["~"]
      ignorePaths: ["~/dev", "~/go", "~/tmp", "~/temp", "~/sdk", "~/Library"]
      ruleType: "file-trigger"
      ruleSettings:
        fileTrigger: "test.file"
        excludePaths:
          - "node_modules"
          - ".venv"

EOF
./heptapod -r ~/test_hepta/rules -v 4 run -d
rm -rf ~/test_hepta
