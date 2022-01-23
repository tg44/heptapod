mkdir -p ~/4test
mkdir -p ~/Abc
mkdir -p ~/Code
touch ~/4test/test.file
touch ~/Abc/test.file
touch ~/Code/test.file
touch ~/Code/.venv
mkdir -p ~/Abc/.venv
mkdir -p ~/test_hepta/rules
cat << EOF > ~/test_hepta/rules/test-rule.yaml
name: "test-node"
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
rm -rf ~/4test
rm -rf ~/Abc
rm -rf ~/Code
