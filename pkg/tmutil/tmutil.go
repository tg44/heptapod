package tmutil

import (
	"bufio"
	"bytes"
	"github.com/tg44/heptapod/pkg/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func AddPathsToTM(paths []string, logPath string, bufferSize int, verbose bool) {
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	currentTime := time.Now()
	defer utils.TimeTrack(currentTime, "tmutil run", verbose)

	logFile, err := os.Create(filepath.Join(logPath, currentTime.Format("2006-01-02_15:04:05")+".log"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = logFile.Close()
	}()

	check := make(chan string, bufferSize)
	add := make(chan string, bufferSize)
	finished := make(chan bool)

	added := 0

	go func() {
		for j := range check {
			res, _ := checkPath(j)
			if res {
				add <- j
			}
		}
		close(add)
	}()

	go func() {
		for j := range add {
			addPath(j, logFile)
			added += 1
		}
		finished <- true
		close(finished)
	}()

	for _, p := range paths {
		check <- p
	}
	close(check)
	<-finished
	if verbose {
		log.Printf("added %d", added)
	}
}

func checkPath(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	cmd := exec.Command("tmutil", "isexcluded", path)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}

	s := out.String()
	return !strings.Contains(s, "[Excluded]"), nil
}

func addPath(path string, logfile *os.File) {
	cmd := exec.Command("tmutil", "addexclusion", path)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	_, _ = logfile.WriteString(path + "\r\n")
}

func removePath(path string) {
	cmd := exec.Command("tmutil", "removeexclusion", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		log.Println(out.String())
		log.Println(outErr.String())
		log.Fatal(err)
	}
}

func RemoveAllFromLogs(logPath string, bufferSize int, verbose bool) {
	files, err := ioutil.ReadDir(logPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		RemoveFileFromLogs(logPath, file.Name(), bufferSize, verbose)
	}
}
func RemoveFileFromLogs(logPath string, fileName string, bufferSize int, verbose bool) {
	file, err := os.Open(filepath.Join(logPath, fileName))
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	RemovePathsFromTM(lines, bufferSize, verbose)
}
func RemovePathsFromTM(paths []string, bufferSize int, verbose bool) {
	defer utils.TimeTrack(time.Now(), "tmutil run", verbose)

	check := make(chan string, bufferSize)
	remove := make(chan string, bufferSize)
	finished := make(chan bool)

	removed := 0

	go func() {
		for j := range check {
			res, err := checkPath(j)
			if !res && err == nil {
				remove <- j
			}
		}
		close(remove)
	}()

	go func() {
		for j := range remove {
			removePath(j)
			removed += 1
		}
		finished <- true
		close(finished)
	}()

	for _, p := range paths {
		check <- p
	}
	close(check)
	<-finished
	if verbose {
		log.Printf("removed %d", removed)
	}
}

func GetExcludeList() string {
	cmd := exec.Command("mdfind", "com_apple_backup_excludeItem = 'com.apple.backupd'")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}

	return out.String()
}
