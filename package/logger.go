package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/PhantomXD-nepal/goauthtemplate/internal/config"
)

type Level string

const (
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
	DEBUG Level = "DEBUG"
)

// ANSI colors
const (
	reset  = "\033[0m"
	gray   = "\033[38;5;245m"
	green  = "\033[38;5;42m"
	yellow = "\033[38;5;214m"
	red    = "\033[38;5;196m"
	blue   = "\033[38;5;39m"
	pink   = "\033[38;5;205m"
)

func Log(level Level, msg string) {
	if config.Envs.Environment != "dev" {
		return
	}

	now := time.Now().Format("15:04:05")
	color, emoji := style(level)

	fmt.Printf(
		"%s%s %s[%s] %s%s%s\n",
		gray, now,
		color, level,
		emoji, msg,
		reset,
	)
}

func Info(msg string)  { Log(INFO, msg) }
func Warn(msg string)  { Log(WARN, msg) }
func Error(msg string) { Log(ERROR, msg) }
func Debug(msg string) { Log(DEBUG, msg) }

func style(level Level) (string, string) {
	switch level {
	case INFO:
		return green, "✨ "
	case WARN:
		return yellow, "⚠️  "
	case ERROR:
		return red, "🔥 "
	case DEBUG:
		return blue, "🧠 "
	default:
		return pink, "🐾 "
	}
}

func Mascot() {
	if config.Envs.Environment != "dev" {
		return
	}

	ascii := `
     ∧＿∧
    ( •‿• )
  ╭─(づ🚀 )─╮
  │ Go Logger │
  ╰──────────╯
`
	fmt.Fprintln(os.Stdout, pink+ascii+reset)
}
