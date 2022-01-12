package utils

import (
	"log"
	"time"
)

func TimeTrack(start time.Time, name string, verbose int) {
	elapsed := time.Since(start)
	if verbose > 0 {
		log.Printf("%s took %s", name, elapsed)
	}
}
