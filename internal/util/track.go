package util

import (
	"time"

	log "github.com/inconshreveable/log15"
)

// Track function execution time
// e.g: defer util.Duration(util.Track("foo"))
// the above goes at top of function

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	log.Info("OP Completed execution", "target", msg, "time", time.Since(start))
}
