# heptapod

This is a command line application to manage and fine-tune
[Time Machine](https://support.apple.com/en-us/HT201250) exclude paths.

## This repository is a WIP! The advertised functionality is nonexistent yet :(

### Repo state
done:
 - architecture the protocol/configs
 - most of asimovs features are ported
 - example rules added

todos:
 - command line interface
 - option to dryrun, show inner states, write to file
 - handle global deps (m2, ivy, nvm, npm) 
 - support tmignore functionality
 - support tmignore like funcionality with dockerignore
 - regexp pattern
 - brew package
 - hook to rerun periodically/before tmbackups
 - port asimov's issues (spotify, spotlight)
 - docker support (at least tell if docker vm is persisted or not)
   - this is kinda easy with tmutil.GetExcludeList() 
 - android vms?
   - this should be also easy
 - preemptive search that tells you how many files with which size will be excluded/included
   - nice to have, we can tell the sizes, but counting files need to actually count the files which could be slow AH 
 - speedtest it
 - purge option
   - we should write down what paths excluded by us, and include them back
 - modify the backup intervals and frequencies
   - not sure it can be done with tmutil, but there are applications for this 
 - probably adding `.` and `..` will break the execution shortcuts and need to fix them
 - wildcards like `*.sh` not working right now, do we want to make them work?

### Notes for migrating to a new machine
`xcode-select --install` may be needed after a migration

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
   - `ignore-files`
 - `settings` other type setting see below

#### Ignore (not yet implemented)
Parses the `.gitignore` or `.dockerignore` files, and excludes its contents.
 - `forceAddPaths` are files that we add even if they would be otherwise excluded
   - like `.gitignore` ignores `.env` files but we forcefully want to add them back 
 - `fileName`
   - `.gitignore` or `.dockerignore`

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
