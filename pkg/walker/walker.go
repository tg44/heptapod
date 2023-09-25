package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tg44/heptapod/pkg/utils"
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
	initialJobCount := len(jobs)
	if initialJobCount == 0 {
		return []string{}
	}

	results := make(chan []string, bufferSize)
	jobQueue := make(chan WalkJob, bufferSize)
	spawn := make(chan WalkJob, bufferSize)

	for i, job := range jobs {
		if (i + 1) <= par {
			spawn <- job
		} else {
			jobQueue <- job
		}
	}

	globalResult := []string{}
	expectedResultCount := initialJobCount
	maxId := 0
	for expectedResultCount > 0 {
		select {
		case singleResult := <-results:
			globalResult = append(globalResult, singleResult...)
			expectedResultCount -= 1
			if len(jobQueue) > 0 {
				// check then act is safe since only this one thread received from jobQueue
				spawn <- <-jobQueue
			}
		case j := <-spawn:
			go walk(maxId, j.Rootpath, j.Walkers, j.AlreadyFiltered, results, jobQueue, verbose)
			maxId += 1
		}
	}
	return globalResult
}

func walk(runnerId int, rootpath string, walkers []Walker, alreadyFiltered []string, results chan []string, jobQueue chan WalkJob, verbose int) {
	defer utils.TimeTrack(time.Now(), fmt.Sprintf("(runner-%d) walk on %s", runnerId, rootpath), verbose)
	hasNext := true
	var res *utils.List = nil
	path, err := utils.FixupPathsToHandleHome(rootpath)
	if err != nil {
		results <- []string{}
		return
	}
	if len(walkers) == 0 {
		results <- []string{}
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
								jobQueue <- WalkJob{f, keeps, alreadyFiltered}
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
	results <- res.ToArray()
}
