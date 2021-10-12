package utils

import (
	"log"
	"time"
)

func TimeTrack(start time.Time, name string, verbose bool) {
	elapsed := time.Since(start)
	if verbose {
		log.Printf("%s took %s", name, elapsed)
	}
}
