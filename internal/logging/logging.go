package logging

import (
	"log"
	"os"
)

// Logger is a thin wrapper around the standard log.Logger to allow future upgrades.
type Logger struct {
	level string
	base  *log.Logger
}

// NewLogger builds a logger; level is accepted for future compatibility.
func NewLogger(level string) *Logger {
	return &Logger{
		level: level,
		base:  log.New(os.Stdout, "magnet2torrent: ", log.LstdFlags),
	}
}

func (l *Logger) Infof(format string, args ...any) {
	l.base.Printf("[INFO] "+format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.base.Printf("[WARN] "+format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.base.Printf("[ERROR] "+format, args...)
}
