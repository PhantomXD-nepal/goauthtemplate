package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/PhantomXD-nepal/goauthtemplate/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("Missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func UuidToBytes(id uuid.UUID) ([]byte, error) {
	return id.MarshalBinary()
}

func BytesToUUID(b []byte) (uuid.UUID, error) {
	return uuid.FromBytes(b)
}
func PrintRoutes(r *mux.Router) {
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		fmt.Printf("ROUTE %-30s METHODS %v\n", path, methods)
		return nil
	})
}

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
