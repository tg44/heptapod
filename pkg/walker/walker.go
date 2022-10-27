package walker

import (
	"fmt"
	"github.com/tg44/heptapod/pkg/utils"
	"os"
	"path/filepath"
	"time"
)

// first ret is the excluding paths we found
// second ret is the global ignore paths if we don't want other runners to go into it
// third ret is the local ignore path if we don't want to go into it next time
// last param is the name of the rule for verbose
type Walker = func(path string, files []os.DirEntry) ([]string, []string, []string, string)

type WalkerIgnores struct {
	w            Walker
	localIgnores []string
}

type WalkJob struct {
	Rootpath        string
	Walkers         []Walker
	AlreadyFiltered []string
}

func Run(jobs []WalkJob, par int, bufferSize int, verbose int) []string {
	defer utils.TimeTrack(time.Now(), "walker run", verbose)
	if len(jobs) == 0 {
		return []string{}
	}
	spawn := make(chan WalkJob, bufferSize)
	start := make(chan bool, bufferSize)
	end := make(chan []string)

	for _, j := range jobs {
		spawn <- j
	}

	for w := 1; w <= par; w++ {
		go worker(w, spawn, start, end, verbose)
	}

	res := []string{}
	//wait loop
	<-start
	for e := range end {
		more := false
		select {
		case <-start:
			more = true
		default:
			more = false
		}
		if !more {
			close(start)
			close(end)
			close(spawn)
		}
		res = append(res, e...)
	}
	return res
}

func worker(id int, spawn chan WalkJob, start chan bool, end chan []string, verbose int) {
	for j := range spawn {
		start <- true
		walk(id, j.Rootpath, j.Walkers, j.AlreadyFiltered, spawn, start, end, verbose)
	}
}

func walk(runnerId int, rootpath string, walkers []Walker, alreadyFiltered []string, spawn chan WalkJob, start chan bool, end chan []string, verbose int) {
	defer utils.TimeTrack(time.Now(), fmt.Sprintf("(runner-%d) walk on %s", runnerId, rootpath), verbose)
	hasNext := true
	var res *utils.List = nil
	path, err := utils.FixupPathsToHandleHome(rootpath)
	if err != nil {
		end <- []string{}
		return
	}
	if len(walkers) == 0 {
		end <- []string{}
		return
	}
	l := &utils.List{path, nil}
	for hasNext {
		next := l.Next
		files, err := os.ReadDir(l.Data)
		if err != nil {
			if verbose > 0 {
				fmt.Println("!!! There was an error: ", err.Error())
			}
			if next == nil {
				hasNext = false
			} else {
				l = next
			}
			continue
		}

		ignores := []WalkerIgnores{}
		globalIgnores := []string{}
		for _, walker := range walkers {
			r, i1, i2, name := walker(l.Data, files)
			if verbose > 1 && len(r) > 0 {
				fmt.Println("-----")
				fmt.Println(name)
				fmt.Println(l.Data)
				for _, v := range r {
					fmt.Println("\t", v)
				}
				if verbose > 2 {
					fmt.Println("-localignore")
					for _, v := range i1 {
						fmt.Println("\t", v)
					}
					if verbose > 3 {
						fmt.Println("-globalignore")
						for _, v := range i2 {
							fmt.Println("\t", v)
						}
						fmt.Println("-alreadyFiltered")
						for _, v := range alreadyFiltered {
							fmt.Println("\t", v)
						}
					}
				}
			}
			res = res.Prepend(r)
			ignores = append(ignores, WalkerIgnores{walker, i2})
			globalIgnores = append(globalIgnores, i1...)
		}

		for _, file := range files {
			if file.IsDir() {
				f := filepath.Join(l.Data, file.Name())
				if !res.Contains(f) {
					if !utils.ContainsSA(alreadyFiltered, f) {
						if !utils.ContainsSA(globalIgnores, f) {
							keeps := []Walker{}
							for _, wi := range ignores {
								if !utils.ContainsSA(wi.localIgnores, f) {
									keeps = append(keeps, wi.w)
								}
							}
							if len(keeps) == len(walkers) {
								//if no ignore happened we go further with this runner
								next = next.AddAsHead(f)
							} else if len(keeps) > 0 {
								//if some ignore happened we ignore the path, and spawn a new job with the nonignorant walkers
								spawn <- WalkJob{f, keeps, alreadyFiltered}
							} else {
								//if keeps is empty we let the next be nil
							}
						}
					}
				}
			}
		}

		if next == nil {
			hasNext = false
		} else {
			l = next
		}
	}
	end <- res.ToArray()
}
