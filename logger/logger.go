package logger

import (
	"log"
	"time"

	"github.com/ve-interactive/influx-protector/rules"
)

// Logger options
type Logger struct {
	Verbose       bool
	SlowQueryTime int64
}

// NewLogger creates and instance of the logger
func NewLogger(verbose bool, slowQueryTime int64) *Logger {
	return &Logger{
		Verbose:       verbose,
		SlowQueryTime: slowQueryTime,
	}
}

// Query logs the query info
func (l *Logger) Query(start time.Time, query string, options *rules.Options) {
	if l.Verbose {
		log.Printf("[QUERY] %s", query)
	}

	elapsed := int64(time.Since(start).Seconds() * 1000)
	if elapsed > options.SlowQuery {
		log.Printf("[SLOWQUERY] '%s' took %dms", query, elapsed)
	}
}

// Error logs the error and the query that generated it
func (l *Logger) Error(query string, err error) {
	log.Printf("[ERROR] %s ('%s')", err, query)
}

// Info logs info messages
func (l *Logger) Info(str string, v interface{}) {
	log.Printf("[INFO] "+str, v)
}
