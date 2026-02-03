package utils

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
	cyan   = "\033[38;5;51m"
	purple = "\033[38;5;141m"
)

func Log(level Level, msg string) {
	if config.Envs.Environment != "dev" {
		return
	}

	now := time.Now().Format("15:04:05")
	color, emoji := style(level)

	// Cute box drawing
	fmt.Printf(
		"%s╭─ %s %s%s%s %s─╮%s\n",
		gray, now, color, emoji, level, gray, reset,
	)
	fmt.Printf(
		"%s│%s  %s%s\n",
		gray, reset, msg, reset,
	)
	fmt.Printf(
		"%s╰─────────────────────────────╯%s\n",
		gray, reset,
	)
}

func Info(msg string) {
	Log(INFO, msg)
}

func Warn(msg string) {
	Log(WARN, msg)
}

func Error(msg string) {
	Log(ERROR, msg)
}

func Debug(msg string) {
	Log(DEBUG, msg)
}

func style(level Level) (string, string) {
	switch level {
	case INFO:
		return green, "✨"
	case WARN:
		return yellow, "⚠️ "
	case ERROR:
		return red, "🔥"
	case DEBUG:
		return blue, "🧠"
	default:
		return pink, "🐾"
	}
}

func Mascot() {
	if config.Envs.Environment != "dev" {
		return
	}

	// 	ascii := `
	//     ╔═══════════════════════════════════════╗
	//     ║                                       ║
	//     ║         /\_/\                        ║
	//     ║        ( o.o )    GoAuth Template     ║
	//     ║         > ^ <                         ║
	//     ║        /|   |\                       ║
	//     ║       (_|   |_)                      ║
	//     ║                                       ║
	//     ║    🌟 Ready to authenticate! 🌟      ║
	//     ║                                       ║
	//     ╚═══════════════════════════════════════╝
	// `

	// Gradient-like effect using multiple colors
	lines := []string{
		"    ╔═══════════════════════════════════════╗",
		"    ║                                       ║",
		"    ║         /\\_/\\                        ║",
		"    ║        ( o.o )    GoAuth Template     ║",
		"    ║         > ^ <                         ║",
		"    ║        /|   |\\                       ║",
		"    ║       (_|   |_)                      ║",
		"    ║                                       ║",
		"    ║    🌟 Server started without 🌟      ║",
		"    ║       any compilation errors!         ║",
		"    ╚═══════════════════════════════════════╝",
	}

	colors := []string{pink, pink, cyan, cyan, purple, purple, pink, cyan, purple, pink, pink}

	for i, line := range lines {
		fmt.Fprintln(os.Stdout, colors[i]+line+reset)
	}
	fmt.Println()
}
