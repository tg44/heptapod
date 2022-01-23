# heptapod

This is a command line application to manage and fine-tune
[Time Machine](https://support.apple.com/en-us/HT201250) exclude paths.

<details>
  <summary>Why it is named after creatures from Arrival?</summary>
Heptapods are extraterrestrial species from the movie Arrival.
They are special because they have non-linear time perspective.
Their written language (Heptapod B) is basically describes 
the future and the past in the same time. Hence the name of the tool.
</details>

### Install

```sh
brew tap tg44/heptapod
```
```sh
brew install heptapod
```

### Usage

To print help;
```
heptapod -h
heptapod <action> -h
```

Move (and use) the default ruleset;
```
mkdir -p ~/.heptapod/rules
cp -R $(brew --prefix heptapod)/rules ~/.heptapod
```

Lists all the rules (you get 4 tables, enabled, disabled, parseable but unrunable and unparsable).
```
heptapod rules ls -a
```

Enable/disable rules (by name), or add/remove ignore folders;
```
heptapod rules disable node bower
heptapod rules enable bower
heptapod rules ignoreAdd ~/.Trash ~/.yarn/cache ~/Library
heptapod rules ignoreRemove ~/.Trash
```

List all the currently excluded TM paths;
```
heptapod tm ls
```

Dryrun the current rules, in verbose mode, also log speed and debug informations. (Potentially list nonexistent directories and files!)
```
heptapod -v 1 run -d
```

To run the current rules, and add the dirs to the TM exclude list. Also writes exclude logs to `~/.heptapod/logs` (or the given `--logDir dir`) for easier revert.
```
heptapod run
```

To revert all the previously added paths from the run-exclude-logs. (`prune -h` could tell you the other useful revert options).
```
heptapod prune -a
```

### Dictionary
 - search path - a path that we want to process
 - ignore path - path that we don't process (further)
 - exclude path - path that we don't want in our TM saves
 - include path - path that we want in our TM saves
 - ignore rule - a rule that parse git/docker ignore file format

### Notes from TM migrating to a new machine
When you try to migrate your TM state to a new machine
`xcode-select --install` may be needed. Somehow this is 
sometimes not migrating as you thought it will.

### Notes from TM in general
There are two ways to exclude a dir from backups;
- exclude by path (`tmutil addexclusion -p`)
  - needs sudo
  - can be read back with `defaults read /Library/Preferences/com.apple.TimeMachine.plist SkipPaths`
  - appears in TM preferences / options
- exclude by flag (`tmutil addexclusion`)
  - keeps the given flag when moved
  - can not be reliably list them (`mdfind com_apple_backup_excludeItem = 'com.apple.backupd'` is a close call, but some folders are excluded by mdfind too)

This tool excludes by flag! You can check any folder manually with `tmutil isexcluded`. If you delete a folder, it will be deleted with its flag. You don't need to clean up ever.
Also, you can only exclude nonexcluded files with tmutil, so we only add them if they are exists and if they are not already added.

### Rules
Every rule has a searchPaths, ignorePaths.
 - `name` for categorization
 - `enabled` for easier enable/disable
 - `searchPaths` are the root of the rule search like `~`
 - `ignorePaths` are subpaths that we want to ignore to make the run quicker
   - like we want to parse dirs under `~`, but not `~/Downloads`
 - `type` is the ruletype
   - `file-trigger`
   - `regexp`
   - `ignore-file`
 - `settings` other type setting see below

#### Ignore file (not yet implemented)
Parses the `.gitignore` or `.dockerignore` files, and excludes its contents.
 - `forceIncludePaths` are files that we add even if they would be otherwise excluded
   - like `.gitignore` ignores `.env` files but we forcefully want to add them back 
 - `fileName`
   - `.gitignore` or `.dockerignore`
 - devdocs:
   - the problem with this, that you need to go into the whole subtree to do the parsings

#### Regexp (not yet implemented)
Ignores all files/folders with the given regexp. This can be slow!
 - `regexp`

#### File trigger
Ignores files/dirs based on other files existence, made for easy language dep ignores.
 - `fileTrigger` like `package.json` or `.git`
 - `excludePaths` like `node-modules` or `.`

### Credits
 - [asimov](https://github.com/stevegrunwell/asimov)
 - [tmignore](https://github.com/samuelmeuli/tmignore)
 - [various stack exchange responses](https://superuser.com/questions/1161038/exclude-folders-by-regex-from-time-machine-backup)
 - [tmutil](https://ss64.com/osx/tmutil.html)

### Contribution
If you are interested in this repo, star it, and write an issue, and we can talk about future ideas there!


### Repo/developement state
done:
- architecture the protocol/configs
- most of asimov's features are ported
- example rules added
- command line interface
- option to dryrun, show inner states, write to file
- purge option
   - we should write down what paths excluded by us, and include them back
- brew package
- ghactions
- rule manage commands (list/enable/disable/ignoreAdd/ignoreRemove)

todos:
- handle global deps (m2, ivy, nvm, npm)
- support tmignore functionality
- support tmignore like funcionality with dockerignore
- regexp pattern
- hook to rerun periodically/before tmbackups
- port asimov's issues (spotify, spotlight)
- docker support (at least tell if docker vm is persisted or not)
   - this is kinda easy with tmutil.GetExcludeList()
- android vms?
   - this should be also easy
- preemptive search that tells you how many files with which size will be excluded/included
   - nice to have, we can tell the sizes, but counting files need to actually count the files which could be slow AH
- speedtest it
- modify the backup intervals and frequencies
   - not sure it can be done with tmutil, but there are applications for this
- probably adding `.` and `..` will break the execution shortcuts and need to fix them
- wildcards like `*.sh` not working right now, do we want to make them work?
