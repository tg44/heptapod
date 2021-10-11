package tmutil

import (
	"bytes"
	"github.com/tg44/heptapod/pkg/utils"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func AddPathsToTM(paths []string, bufferSize int) {
	defer utils.TimeTrack(time.Now(), "tmutil run")
	check := make(chan string, bufferSize)
	add := make(chan string, bufferSize)
	finished := make(chan bool)

	added := 0

	go func() {
		for j := range check {
			res := checkPath(j)
			if res {
				add <- j
			}
		}
		close(add)
	}()

	go func() {
		for j := range add {
			addPath(j)
			added += 1
		}
		finished <- true
		close(finished)
	}()

	for _, p := range paths {
		check <- p
	}
	close(check)
	<- finished
	log.Printf("added %d", added)
}

func checkPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	cmd := exec.Command("tmutil", "isexcluded", path)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("?")
		log.Println(path)
		log.Println(err)
		return false
	}

	s := out.String()
	return !strings.Contains(s, "[Excluded]")
}

func addPath(path string) {
	cmd := exec.Command("tmutil", "addexclusion", path)
	err := cmd.Run()
	if err != nil {
		log.Println("!")
		log.Println(err)
	}
}

func GetExcludeList() string {
	cmd := exec.Command("mdfind",  "com_apple_backup_excludeItem = 'com.apple.backupd'")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}

	return out.String()
}
