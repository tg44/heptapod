package walker

import (
	"fmt"
	"github.com/tg44/heptapod/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

//first ret is the excluding paths we found
//second ret is the global ignore paths if we don't want other runners to go into it
//third ret is the local ignore path if we don't want to go into it next time
type Walker = func(path string, files []os.FileInfo) ([]string, []string, []string)

type WalkerIgnores struct {
	w            Walker
	localIgnores []string
}

type WalkJob struct {
	Rootpath        string
	Walkers         []Walker
	AlreadyFiltered []string
}

func Run(jobs []WalkJob, par int, bufferSize int, verbose bool) []string {
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

func worker(id int, spawn chan WalkJob, start chan bool, end chan []string, verbose bool) {
	for j := range spawn {
		walk(id, j.Rootpath, j.Walkers, j.AlreadyFiltered, spawn, start, end, verbose)
	}
}

func walk(runnerId int, rootpath string, walkers []Walker, alreadyFiltered []string, spawn chan WalkJob, start chan bool, end chan []string, verbose bool) {
	defer utils.TimeTrack(time.Now(), fmt.Sprintf("(runner-%d) walk on %s", runnerId, rootpath), verbose)
	start <- true
	hasNext := true
	var res *utils.List = nil
	path, err := utils.FixupPathsToHandleHome(rootpath)
	if err != nil {
		end <- []string{}
	}
	l := &utils.List{path, nil}
	for hasNext {
		next := l.Next
		files, err := ioutil.ReadDir(l.Data)
		if err != nil {
			if next == nil {
				hasNext = false
			} else {
				l = next
			}
			break
		}

		ignores := []WalkerIgnores{}
		globalIgnores := []string{}
		for _, walker := range walkers {
			r, i1, i2 := walker(l.Data, files)
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
								next = next.AddAsHead(filepath.Join(l.Data, file.Name()))
							} else {
								//if some ignore happened we ignore the path, and spawn a new job with the nonignorant walkers
								spawn <- WalkJob{filepath.Join(l.Data, file.Name()), keeps, alreadyFiltered}
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
