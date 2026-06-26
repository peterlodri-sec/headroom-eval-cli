package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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
	hidden  = lipgloss.NewStyle().Foreground(lipgloss.Color("#0a0a12")) // steganographic — same as bg

	// Genesis seal — embedded at build time via -ldflags
	Version   = "dev"
	Seal      = "genesis"
	BuildTime = "unknown"
)

// Hidden ASCII art — rendered in background color (steganographic)
// Visible only when terminal bg mismatches or on copy-paste
var stegoArt = strings.Repeat(" ", 4) + hidden.Render(strings.Join([]string{
	"░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
	"░░  loopkit · kompress-v8 · iclr 2027  ░░",
	"░░  the loop shipped. the paradox is proven.  ░░",
	"░░  label quality is the bottleneck.  ░░",
	"░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
}, "\n"+strings.Repeat(" ", 4)))

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
	resp, err := http.Get("https://huggingface.co/api/spaces/PeetPedro/headroom-eval")
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
		case "s":
			m.showStego = !m.showStego // toggle steganography visibility
		}
	}
	return m, nil
}

func openURL(url string) {
	fmt.Fprintf(os.Stderr, "\n→ %s\n", url)
}

// genesisHash returns a deterministic seal of the binary
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
		return "Window too small — please resize"
	}

	var b strings.Builder

	// Animated ASCII dodecahedron
	frames := []string{
		"   ╭──────╮    ╭──────╮\n  ╱ ╲    ╱ ╲  ╱ ╲    ╱ ╲\n ╱   ╲  ╱   ╲╱   ╲  ╱   ╲\n╱     ╲╱     ╲     ╲╱     ╲\n╲     ╱╲     ╱     ╱╲     ╱\n ╲   ╱  ╲   ╱╲   ╱  ╲   ╱\n  ╲ ╱    ╲ ╱  ╲ ╱    ╲ ╱\n   ╰──────╯    ╰──────╯",
		"   ╭──────╮    ╭──────╮\n  ╱ ╲    ╱ ╲  ╱    ╲  ╱ ╲\n ╱   ╲  ╱   ╲╱      ╲╱   ╲\n╱     ╲╱     ╲        ╲     ╲\n╲     ╱╲     ╱        ╱     ╱\n ╲   ╱  ╲   ╱╲      ╱╲   ╱\n  ╲ ╱    ╲ ╱  ╲    ╱  ╲ ╱\n   ╰──────╯    ╰──────╯",
	}
	frame := frames[m.frame%len(frames)]

	b.WriteString(title.Render("🐋 headroom-eval"))
	b.WriteString(subtle.Render(fmt.Sprintf("  v%s  seal:%s", Version, genesisHash())))
	b.WriteString("\n")
	b.WriteString(subtle.Render(frame))
	b.WriteString("\n")

	// Steganographic layer — visible only with 's' key
	if m.showStego {
		b.WriteString(stegoArt)
		b.WriteString("\n")
	}

	// Status
	if m.err != "" {
		b.WriteString(failure.Render(fmt.Sprintf("  ✗ %s", m.err)))
	} else if m.ready {
		icon, style := "✅", success
		if strings.Contains(m.spaceStatus, "BUILD") || strings.Contains(m.spaceStatus, "build") {
			icon, style = "🔨", subtle
		} else if m.spaceStatus != "running" && m.spaceStatus != "RUNNING" {
			icon, style = "⚠️", failure
		}
		b.WriteString(style.Render(fmt.Sprintf("  %s Space: %s", icon, m.spaceStatus)))
	} else {
		b.WriteString(subtle.Render("  ⏳ connecting..."))
	}
	b.WriteString("\n\n")

	k := func(s string) string { return lipgloss.NewStyle().Foreground(accent).Bold(true).Render(fmt.Sprintf("[%s]", s)) }
	b.WriteString(subtle.Render("  keys:"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s refresh  %s Space  %s paper  %s loopkit  %s stego  %s quit\n", k("r"), k("o"), k("p"), k("g"), k("s"), k("q")))
	b.WriteString("\n")
	b.WriteString(subtle.Render("  headroom eval · hill climbing loop · iclr 2027"))

	return b.String()
}

func main() {
	// CLI flags for non-TUI mode
	sealFlag := flag.Bool("seal", false, "print genesis seal and exit")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

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
		fmt.Fprintf(os.Stderr, "whale: %v\n", err)
		os.Exit(1)
	}
}
