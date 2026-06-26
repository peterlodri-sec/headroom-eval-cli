package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	dim    = "\033[2m"
	reset  = "\033[0m"
)

func main() {
	token := flag.String("token", os.Getenv("HF_TOKEN"), "HF API token")
	space := flag.String("space", "PeetPedro/headroom-eval", "HF Space")
	mode := flag.String("mode", "run", "run or build")
	follow := flag.Bool("f", false, "follow (stream)")
	filter := flag.String("filter", "", "filter lines containing this string")
	flag.Parse()

	if *token == "" {
		fmt.Fprintln(os.Stderr, "set HF_TOKEN or use --token")
		os.Exit(1)
	}

	url := fmt.Sprintf("https://huggingface.co/api/spaces/%s/logs/%s", *space, *mode)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+*token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%serror:%s %v\n", red, reset, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 1MB buffer

	for scanner.Scan() {
		line := scanner.Text()

		if *filter != "" && !strings.Contains(line, *filter) {
			continue
		}

		// Colorize
		switch {
		case strings.Contains(line, "ERROR") || strings.Contains(line, "Traceback") || strings.Contains(line, "FAILED"):
			fmt.Println(red + line + reset)
		case strings.Contains(line, "WARNING") || strings.Contains(line, "WARN"):
			fmt.Println(yellow + line + reset)
		case strings.Contains(line, "INFO") || strings.Contains(line, "success") || strings.Contains(line, "complete"):
			fmt.Println(green + line + reset)
		case strings.Contains(line, "===") || strings.Contains(line, "Starting") || strings.Contains(line, "Running"):
			fmt.Println(cyan + line + reset)
		default:
			fmt.Println(dim + line + reset)
		}

		if !*follow {
			break // one-shot if not following
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%serror:%s %v\n", red, reset, err)
	}
}
