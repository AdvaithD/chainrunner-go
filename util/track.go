package util

import (
	"log"
	"time"
)

// Track function execution time
// e.g: defer util.Duration(util.Track("foo"))
// the above goes at top of function


func Track(msg string) (string, time.Time) {
    return msg, time.Now()
}

func Duration(msg string, start time.Time) {
    log.Printf("%v ---- %v\n", msg, time.Since(start))
}