package main

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles du TUI.
var (
	titleStyle   = lipgloss.NewStyle().Bold(true)
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	okStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	missingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

type viewKind int

const (
	viewList viewKind = iota
	viewDetail
)

type model struct {
	configs []Config
	cursor  int
	view    viewKind

	loading bool
	busy    bool
	status  string
	err     error
}

// messages.
type loadedMsg struct {
	configs []Config
	err     error
}

type finishedMsg struct{ err error }

func newModel() model {
	return model{loading: true}
}

func runTUI() error {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func discoverCmd() tea.Cmd {
	return func() tea.Msg {
		configs, err := discover()
		return loadedMsg{configs, err}
	}
}

func (m model) Init() tea.Cmd {
	return discoverCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case loadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.configs = msg.configs
			m.clampCursor()
		}
		return m, nil

	case finishedMsg:
		m.busy = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.status = ""
		}
		m.loading = true
		return m, discoverCmd()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Toujours pouvoir quitter.
	if key == "q" || key == "ctrl+c" {
		return m, tea.Quit
	}
	if m.loading || m.busy {
		return m, nil
	}

	switch key {
	case "up", "k":
		if m.view == viewList {
			m.moveCursor(-1)
		}
	case "down", "j":
		if m.view == viewList {
			m.moveCursor(1)
		}
	case "enter", "l":
		if m.view == viewList && len(m.visible()) > 0 {
			m.view = viewDetail
		}
	case "esc", "h":
		if m.view == viewDetail {
			m.view = viewList
		}
	case "r":
		m.loading = true
		m.status = ""
		return m, discoverCmd()
	case "i", "u", "a":
		return m.action(key)
	}
	return m, nil
}

// action prépare et lance l'installation/mise à jour (courante en vue détails,
// globale en vue liste).
func (m model) action(key string) (tea.Model, tea.Cmd) {
	var steps []step
	var warns []string

	if m.view == viewDetail {
		c := m.current()
		if c == nil {
			return m, nil
		}
		s, err := stepsFor(c)
		if err != nil {
			m.err = err
			return m, nil
		}
		steps = []step{s}
	} else {
		for i := range m.configs {
			c := &m.configs[i]
			if c.Hidden {
				continue
			}
			if key == "i" && c.Installed {
				continue
			}
			if key == "u" && !c.Installed {
				continue
			}
			// En liste, "a" se comporte comme "u" (mettre à jour).
			if key == "a" && !c.Installed {
				continue
			}
			s, err := stepsFor(c)
			if err != nil {
				warns = append(warns, err.Error())
				continue
			}
			steps = append(steps, s)
		}
	}

	if len(steps) == 0 {
		if len(warns) > 0 {
			m.err = fmt.Errorf("%s", strings.Join(warns, "\n"))
		} else {
			m.status = "rien à faire"
		}
		return m, nil
	}

	path, cleanup, err := buildScript(steps)
	if err != nil {
		m.err = err
		return m, nil
	}
	cmd := exec.Command("sh", path)
	m.busy = true
	m.status = "exécution…"
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		cleanup()
		return finishedMsg{err}
	})
}

// ----- navigation -----------------------------------------------------------

// visible renvoie les index des configs affichables (hors HIDDEN).
func (m model) visible() []int {
	var idx []int
	for i := range m.configs {
		if !m.configs[i].Hidden {
			idx = append(idx, i)
		}
	}
	return idx
}

func (m *model) clampCursor() {
	n := len(m.visible())
	if n == 0 {
		m.cursor = 0
		return
	}
	if m.cursor >= n {
		m.cursor = n - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *model) moveCursor(delta int) {
	n := len(m.visible())
	if n == 0 {
		return
	}
	m.cursor = (m.cursor + delta + n) % n
}

func (m model) current() *Config {
	vis := m.visible()
	if m.cursor < 0 || m.cursor >= len(vis) {
		return nil
	}
	return &m.configs[vis[m.cursor]]
}

// ----- rendu ----------------------------------------------------------------

func (m model) View() string {
	switch {
	case m.loading && len(m.configs) == 0:
		return titleStyle.Render("configs") + "\n\n" + "Chargement des configurations…\n"
	case m.view == viewDetail:
		return m.detailView()
	default:
		return m.listView()
	}
}

func (m model) listView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("configs"))
	b.WriteString(dimStyle.Render("  —  " + repoSlug))
	b.WriteString("\n\n")

	vis := m.visible()
	if len(vis) == 0 {
		b.WriteString(dimStyle.Render("aucune configuration trouvée") + "\n")
	}
	for i, ci := range vis {
		c := m.configs[ci]
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
		}
		mark, word := statusMark(&c)
		fmt.Fprintf(&b, "%s%s %-10s %-22s %s\n",
			cursor, mark, c.label(), c.TargetPath, dimStyle.Render(word+"  "+c.Desc))
	}

	b.WriteString("\n")
	b.WriteString(m.footer("↑/↓ naviguer · entrée détails · i installer les manquants · u tout mettre à jour · r actualiser · q quitter"))
	return b.String()
}

func (m model) detailView() string {
	c := m.current()
	if c == nil {
		return m.listView()
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render(c.label()))
	if c.Hidden {
		b.WriteString(dimStyle.Render(" (masqué)"))
	}
	b.WriteString("\n")
	if c.Desc != "" {
		b.WriteString(c.Desc + "\n")
	}
	b.WriteString("\n")

	mark, word := statusMark(c)
	fmt.Fprintf(&b, "%s : %s  %s\n", titleStyle.Render("statut"), word, mark)
	fmt.Fprintf(&b, "%s : %s\n", titleStyle.Render("cible"), c.TargetPath)
	if c.GitManaged {
		fmt.Fprintf(&b, "%s : git, branche %s\n", titleStyle.Render("gestion"), c.GitBranch)
	}
	if c.Desync {
		b.WriteString(warnStyle.Render("cible présente mais non gérée par git (désync) — intervention manuelle") + "\n")
	}
	if c.LastCommit != "" {
		fmt.Fprintf(&b, "%s : %s\n", titleStyle.Render("dernier commit"), c.LastCommit)
	}
	if len(c.Prereq) > 0 {
		b.WriteString(titleStyle.Render("pré-requis") + " :")
		for _, p := range c.Prereq {
			mark := okStyle.Render("✓")
			if !c.Prereqs[p] {
				mark = missingStyle.Render("✗")
			}
			b.WriteString(" " + p + mark)
		}
		b.WriteString("\n")
	}

	action := "installer"
	if c.Installed {
		action = "mettre à jour"
	}
	b.WriteString("\n")
	b.WriteString(m.footer("a/i/u " + action + " · esc retour · q quitter"))
	return b.String()
}

func (m model) footer(help string) string {
	var b strings.Builder
	if m.err != nil {
		b.WriteString(errStyle.Render(m.err.Error()) + "\n")
	}
	if m.status != "" {
		b.WriteString(dimStyle.Render(m.status) + "\n")
	}
	b.WriteString(helpStyle.Render(help))
	return b.String()
}

// statusMark renvoie l'icône et le mot d'état d'une config.
func statusMark(c *Config) (string, string) {
	switch {
	case !c.Installed:
		return missingStyle.Render("✗"), "absent"
	case c.Desync:
		return warnStyle.Render("!"), "désync"
	default:
		return okStyle.Render("✓"), "installé"
	}
}
