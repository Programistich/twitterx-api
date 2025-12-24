package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var debugEnabled bool

func init() {
	debugEnabled = os.Getenv("DEBUG") != ""
}

func formatMessage(level, format string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, msg)
}

// Info logs informational messages (always enabled)
func Info(format string, args ...interface{}) {
	log.Println(formatMessage("INFO", format, args...))
}

// Debug logs debug messages (only when DEBUG env is set)
func Debug(format string, args ...interface{}) {
	if debugEnabled {
		log.Println(formatMessage("DEBUG", format, args...))
	}
}

// Error logs error messages (always enabled)
func Error(format string, args ...interface{}) {
	log.Println(formatMessage("ERROR", format, args...))
}

// Fatal logs error message and exits
func Fatal(format string, args ...interface{}) {
	log.Fatal(formatMessage("ERROR", format, args...))
}

// IsDebugEnabled returns whether debug logging is enabled
func IsDebugEnabled() bool {
	return debugEnabled
}
