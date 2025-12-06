package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Logger is a thin wrapper around the standard log.Logger to allow future upgrades.
type Logger struct {
	level string
	base  *log.Logger
}

// NewLogger builds a logger; level is accepted for future compatibility.
func NewLogger(level string, logPath string) *Logger {
	writers := []io.Writer{os.Stdout}

	if logPath != "" {
		if f, err := openLogFile(logPath); err != nil {
			fmt.Fprintf(os.Stderr, "magnet2torrent: could not open log file %s: %v\n", logPath, err)
		} else {
			writers = append(writers, f)
		}
	}

	return &Logger{
		level: level,
		base:  log.New(io.MultiWriter(writers...), "magnet2torrent: ", log.LstdFlags),
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

func openLogFile(path string) (io.Writer, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir %s: %w", dir, err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644) // #nosec G302
	if err != nil {
		return nil, fmt.Errorf("open log file %s: %w", path, err)
	}
	return f, nil
}
