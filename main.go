package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	accent  = lipgloss.Color("#3ee8c5")
	dim     = lipgloss.Color("#565f89")
	errC    = lipgloss.Color("#f7768e")
	okC     = lipgloss.Color("#9ece6a")
	title   = lipgloss.NewStyle().Foreground(accent).Bold(true)
	subtle  = lipgloss.NewStyle().Foreground(dim)
	success = lipgloss.NewStyle().Foreground(okC)
	failure = lipgloss.NewStyle().Foreground(errC)
	hidden  = lipgloss.NewStyle().Foreground(lipgloss.Color("#0a0a12"))

	Version   = "dev"
	BuildTime = "unknown"
)

// Pure ASCII banner — no Unicode, works everywhere
var asciiBanner = strings.Join([]string{
	"    _    _                _                                   _",
	"   | |  | |              | |                                 | |",
	"   | |__| | ___  __ _  __| |_ __ ___   ___  _ __      ___   | |",
	"   |  __  |/ _ \\/ _` |/ _` | '__/ _ \\ / _ \\| '_ \\    / _ \\  | |",
	"   | |  | |  __/ (_| | (_| | | | (_) | (_) | | | |  |  __/  | |",
	"   |_|  |_|\\___|\\__,_|\\__,_|_|  \\___/ \\___/|_| |_|   \\___|  |_|",
	"",
	"   headroom eval · hill climbing loop · iclr 2027",
	"   github.com/peterlodri-sec · pocoo.vaked.dev",
	"", "",
}, "\n")

var bannerFrames = []string{
	strings.Join([]string{
		"       +--------+        +--------+",
		"      / \\      / \\      / \\      / \\",
		"     /   \\    /   \\    /   \\    /   \\",
		"    /     \\  /     \\  /     \\  /     \\",
		"    \\     /  \\     /  \\     /  \\     /",
		"     \\   /    \\   /    \\   /    \\   /",
		"      \\ /      \\ /      \\ /      \\ /",
		"       +--------+        +--------+",
	}, "\n"),
	strings.Join([]string{
		"       +--------+        +--------+",
		"      / \\      / \\      /        \\/ \\",
		"     /   \\    /   \\    /          /   \\",
		"    /     \\  /     \\  /          /     \\",
		"    \\     /  \\     /  \\          \\     /",
		"     \\   /    \\   /    \\          \\   /",
		"      \\ /      \\ /      \\        / \\ /",
		"       +--------+        +--------+",
	}, "\n"),
}

type model struct {
	spaceStatus string
	ready       bool
	err         string
	width       int
	height      int
	frame       int
	showStego   bool
}

type statusMsg struct {
	status string
	err    error
}

func fetchSpace() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://huggingface.co/api/spaces/PeetPedro/headroom-eval", nil)
	if err != nil {
		return statusMsg{err: err}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return statusMsg{err: err}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	s, _ := data["status"].(string)
	return statusMsg{status: s}
}

func tickCmd() tea.Msg {
	time.Sleep(3 * time.Second)
	return tickMsg{}
}

type tickMsg struct{}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchSpace, tickCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case statusMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
		} else {
			m.spaceStatus = msg.status
			m.ready = true
			m.err = ""
		}
	case tickMsg:
		m.frame++
		return m, tea.Batch(fetchSpace, tickCmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			return m, fetchSpace
		case "o":
			openURL("https://huggingface.co/spaces/PeetPedro/headroom-eval")
		case "p":
			openURL("https://kompress.vaked.dev")
		case "g":
			openURL("https://github.com/peterlodri-sec/loopkit")
		case "d":
			openURL("https://github.com/sponsors/peterlodri-sec")
		case "s":
			m.showStego = !m.showStego
		}
	}
	return m, nil
}

func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "\n-> %s\n", url)
	}
}

func genesisHash() string {
	exe, _ := os.Executable()
	data, err := os.ReadFile(exe)
	if err != nil {
		return "unsealed"
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:8])
}

func (m model) View() string {
	if m.width < 60 {
		return asciiBanner + "\n[window too small — please resize]"
	}

	var b strings.Builder

	// ASCII animated dodecahedron
	b.WriteString(title.Render("headroom-eval"))
	b.WriteString(subtle.Render(fmt.Sprintf("  v%s  seal:%s", Version, genesisHash())))
	b.WriteString("\n")
	b.WriteString(subtle.Render(bannerFrames[m.frame%len(bannerFrames)]))
	b.WriteString("\n")

	// Steganography
	if m.showStego {
		stego := []string{
			"  loopkit . kompress-v8 . iclr 2027",
			"  the loop shipped. the paradox is proven.",
			"  label quality is the bottleneck.",
		}
		for _, s := range stego {
			b.WriteString(hidden.Render(s))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Status
	if m.err != "" {
		b.WriteString(failure.Render(fmt.Sprintf("  x %s", m.err)))
	} else if m.ready {
		icon, style := "+", success
		if strings.Contains(m.spaceStatus, "BUILD") || strings.Contains(m.spaceStatus, "build") {
			icon, style = "#", subtle
		} else if m.spaceStatus != "running" && m.spaceStatus != "RUNNING" {
			icon, style = "!", failure
		}
		b.WriteString(style.Render(fmt.Sprintf("  %s Space: %s", icon, m.spaceStatus)))
	} else {
		b.WriteString(subtle.Render("  ... connecting ..."))
	}
	b.WriteString("\n\n")

	// Keys
	k := func(s string) string { return lipgloss.NewStyle().Foreground(accent).Bold(true).Render(fmt.Sprintf("[%s]", s)) }
	links := []struct{ key, label string }{
		{"r", "refresh"}, {"o", "Space"}, {"p", "paper"}, {"g", "loopkit"},
		{"d", "donate"}, {"s", "stego"}, {"q", "quit"},
	}
	b.WriteString(subtle.Render("  keys: "))
	for i, l := range links {
		b.WriteString(k(l.key))
		b.WriteString(l.label)
		if i < len(links)-1 {
			b.WriteString("  ")
		}
	}

	return b.String()
}

func main() {
	helpFlag := flag.Bool("help", false, "show help")
	sealFlag := flag.Bool("seal", false, "print genesis seal and exit")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *helpFlag {
		fmt.Println("headroom-eval — interactive ASCII TUI")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  headroom-eval              launch TUI")
		fmt.Println("  headroom-eval --seal       print genesis hash")
		fmt.Println("  headroom-eval --version    print version")
		fmt.Println("  headroom-eval --help       this message")
		fmt.Println("")
		fmt.Println("TUI keys:")
		fmt.Println("  r  refresh    o  Space    p  paper    g  loopkit")
		fmt.Println("  d  donate     s  stego    q  quit")
		fmt.Println("")
		fmt.Println("Install: go install github.com/peterlodri-sec/headroom-eval-cli@latest")
		return
	}

	if *sealFlag {
		fmt.Printf("genesis: %s\nversion: %s\nbuilt:   %s\n", genesisHash(), Version, BuildTime)
		return
	}
	if *versionFlag {
		fmt.Printf("headroom-eval %s (%s)\n", Version, BuildTime)
		return
	}

	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "headroom-eval: %v\n", err)
		os.Exit(1)
	}
}
