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

func AddPathsToTM(paths []string, logDir string, bufferSize int, verbose int) int {
	logPath, err := utils.FixupPathsToHandleHome(logDir)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(logPath, os.ModePerm)
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
			res, err := checkPath(j, verbose)
			if res && err == nil {
				add <- j
			}
		}
		close(add)
	}()

	go func() {
		for j := range add {
			err = addPath(j, logFile)
			if err == nil {
				added += 1
			}
		}
		finished <- true
		close(finished)
	}()

	for _, p := range paths {
		check <- p
	}
	close(check)
	<-finished
	if verbose > 0 {
		log.Printf("added %d", added)
	}
	return added
}

func checkPath(path string, verbose int) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	cmd := exec.Command("tmutil", "isexcluded", path)

	var out bytes.Buffer
	cmd.Stdout = &out
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		if verbose > 1 {
			log.Println(out.String())
			log.Println(outErr.String())
		}
		return false, err
	}

	s := out.String()
	return !strings.Contains(s, "[Excluded]"), nil
}

func addPath(path string, logfile *os.File) error {
	cmd := exec.Command("tmutil", "addexclusion", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		log.Println(out.String())
		log.Println(outErr.String())
		return err
	}
	_, _ = logfile.WriteString(path + "\r\n")
	return nil
}

func removePath(path string) error {
	cmd := exec.Command("tmutil", "removeexclusion", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		log.Println(out.String())
		log.Println(outErr.String())
	}
	return err
}

func RemoveAllFromLogs(logPath string, bufferSize int, verbose int) {
	files, err := ioutil.ReadDir(logPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		RemoveFileFromLogs(logPath, file.Name(), bufferSize, verbose)
	}
}
func RemoveFileFromLogs(logPath string, fileName string, bufferSize int, verbose int) {
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
func RemovePathsFromTM(paths []string, bufferSize int, verbose int) {
	defer utils.TimeTrack(time.Now(), "tmutil run", verbose)

	check := make(chan string, bufferSize)
	remove := make(chan string, bufferSize)
	finished := make(chan bool)

	removed := 0

	go func() {
		for j := range check {
			res, err := checkPath(j, verbose)
			if !res && err == nil {
				remove <- j
			}
		}
		close(remove)
	}()

	go func() {
		for j := range remove {
			err := removePath(j)
			if err == nil {
				removed += 1
			}
		}
		finished <- true
		close(finished)
	}()

	for _, p := range paths {
		check <- p
	}
	close(check)
	<-finished
	if verbose > 0 {
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

func GetExcludeListAll() string {
	return "We can't list all of the excluded folders :(\nThe reason is mainly because there are two modes to get a file/folder excluded by TM;\n\n" +
		" - we add a path as-is, 'hardcode' it as you wish, these can be read by `defaults read /Library/Preferences/com.apple.TimeMachine.plist SkipPaths` and these are the ones appear in the TM settings window\n\n" +
		" - we can 'flag' files/folders, and if you move them they will be still flagged, but these are not in any list, so you either traverse every folder and check if it is excluded or not, or you use a cache like mdfind which is not a full cache (some files are excluded from it)\n\n" +
		"So we can either add files by `tmutil addexclusion` and we can't list them all anymore, or we use `tmutil addexclusion -p` which will flood the settings, and more easy to break and also needs sudo. I chose the first option, so we can't reliably list all the excluded folders ATM.\n\n" +
		"If you have any idea, or know a way to do this, contributions (or at least a new issue) welcomed at https://github.com/tg44/heptapod"
}
