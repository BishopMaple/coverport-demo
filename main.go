package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	_ "github.com/BishopMaple/coverport-demo/internal/coverage" // coverport: coverage instrumentation
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/echo", handleEcho)
	mux.HandleFunc("/api/calc", handleCalc)
	mux.HandleFunc("/api/time", handleTime)
	mux.HandleFunc("/api/info", handleInfo)

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>CoverPort Demo</h1>
<p>Endpoints:</p>
<ul>
<li><a href="/health">/health</a></li>
<li><a href="/api/echo?msg=hello">/api/echo?msg=hello</a></li>
<li><a href="/api/calc?op=add&a=2&b=3">/api/calc?op=add&a=2&b=3</a></li>
<li><a href="/api/time">/api/time</a></li>
<li><a href="/api/info">/api/info</a></li>
</ul>`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		http.Error(w, "missing 'msg' parameter", http.StatusBadRequest)
		return
	}
	upper := strings.ToUpper(msg)
	reversed := reverseString(msg)
	writeJSON(w, map[string]string{
		"original": msg,
		"upper":    upper,
		"reversed": reversed,
		"length":   strconv.Itoa(len(msg)),
	})
}

func handleCalc(w http.ResponseWriter, r *http.Request) {
	op := r.URL.Query().Get("op")
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		http.Error(w, "invalid parameter 'a'", http.StatusBadRequest)
		return
	}
	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		http.Error(w, "invalid parameter 'b'", http.StatusBadRequest)
		return
	}

	var result float64
	switch op {
	case "add":
		result = a + b
	case "sub":
		result = a - b
	case "mul":
		result = a * b
	case "div":
		if b == 0 {
			http.Error(w, "division by zero", http.StatusBadRequest)
			return
		}
		result = a / b
	case "pow":
		result = math.Pow(a, b)
	case "mod":
		if b == 0 {
			http.Error(w, "modulo by zero", http.StatusBadRequest)
			return
		}
		result = math.Mod(a, b)
	default:
		http.Error(w, "unknown op: use add, sub, mul, div, pow, mod", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]interface{}{
		"op":     op,
		"a":      a,
		"b":      b,
		"result": result,
	})
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	writeJSON(w, map[string]interface{}{
		"utc":       now.Format(time.RFC3339),
		"unix":      now.Unix(),
		"day":       now.Weekday().String(),
		"week":      weekNumber(now),
		"leap_year": isLeapYear(now.Year()),
	})
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	writeJSON(w, map[string]interface{}{
		"hostname": hostname,
		"pid":      os.Getpid(),
		"go_env":   os.Getenv("GOVERSION"),
		"port":     os.Getenv("PORT"),
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func weekNumber(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
