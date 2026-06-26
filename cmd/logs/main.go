package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	green  string
	red    string
	yellow string
	cyan   string
	dim    string
	reset  string
)

func initColors() {
	// Unix conventions: no color if piped, NO_COLOR set, or TERM=dumb
	color := true
	if os.Getenv("NO_COLOR") != "" {
		color = false
	}
	if os.Getenv("TERM") == "dumb" {
		color = false
	}
	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == 0 {
		color = false // piped — respect Unix conventions
	}

	if color {
		green = "\033[32m"
		red = "\033[31m"
		yellow = "\033[33m"
		cyan = "\033[36m"
		dim = "\033[2m"
	}
	reset = "\033[0m"
}

func main() {
	token := flag.String("token", os.Getenv("HF_TOKEN"), "HF API token ($HF_TOKEN)")
	space := flag.String("space", "PeetPedro/headroom-eval", "HF Space (owner/name)")
	mode := flag.String("mode", "run", "run or build")
	follow := flag.Bool("f", false, "follow (stream)")
	filter := flag.String("filter", "", "filter lines containing STR")
	colorFlag := flag.Bool("color", false, "force color (overrides auto-detect)")
	flag.Parse()

	initColors()
	if *colorFlag {
		green = "\033[32m"
		red = "\033[31m"
		yellow = "\033[33m"
		cyan = "\033[36m"
		dim = "\033[2m"
	}

	if *token == "" {
		fmt.Fprintln(os.Stderr, "headroom-logs: HF_TOKEN not set. Use --token or export HF_TOKEN.")
		os.Exit(1)
	}

	url := fmt.Sprintf("https://huggingface.co/api/spaces/%s/logs/%s", *space, *mode)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "headroom-logs: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+*token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "headroom-logs: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Handle SIGPIPE gracefully — Unix expectation
	signalPIPE()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if *filter != "" && !strings.Contains(line, *filter) {
			continue
		}

		// Colorize only if terminal
		c := ""
		switch {
		case strings.Contains(line, "ERROR") || strings.Contains(line, "Traceback") || strings.Contains(line, "FAILED") || strings.Contains(line, "ModuleNotFoundError"):
			c = red
		case strings.Contains(line, "WARNING") || strings.Contains(line, "WARN"):
			c = yellow
		case strings.Contains(line, "INFO") || strings.Contains(line, "successfully") || strings.Contains(line, "complete"):
			c = green
		case strings.Contains(line, "===") || strings.Contains(line, "Starting") || strings.Contains(line, "Running") || strings.Contains(line, "Building"):
			c = cyan
		default:
			if !*colorFlag {
				c = "" // no dim in pipe mode
			} else {
				c = dim
			}
		}

		if c != "" {
			fmt.Print(c + line + reset + "\n")
		} else {
			fmt.Println(line)
		}

		if !*follow {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		// SIGPIPE is expected when piped to head/tail — don't error
		if !isPipeError(err) {
			fmt.Fprintf(os.Stderr, "headroom-logs: %v\n", err)
		}
	}
}

func signalPIPE() {
	// On SIGPIPE, just exit 0 — Unix convention for pipe chains.
	// This goroutine is intentionally long-lived for the process lifetime.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGPIPE)
	go func() {
		<-ch
		os.Exit(0)
	}()
}

func isPipeError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "broken pipe") || strings.Contains(s, "write /dev/stdout")
}
