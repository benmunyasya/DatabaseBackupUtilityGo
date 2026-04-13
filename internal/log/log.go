// internal/log/log.go
package log

import (
    "fmt"
)

// ANSI color codes
const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Cyan   = "\033[36m"
)

// Info messages (cyan)
func Info(msg string) {
    fmt.Println(Cyan + "[INFO] " + msg + Reset)
}

// Success messages (green)
func Success(msg string) {
    fmt.Println(Green + "[SUCCESS] " + msg + Reset)
}

// Warning messages (yellow)
func Warn(msg string) {
    fmt.Println(Yellow + "[WARN] " + msg + Reset)
}

// Error messages (red)
func Error(msg string) {
    fmt.Println(Red + "[ERROR] " + msg + Reset)
}
