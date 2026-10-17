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
	titleStyle    = lipgloss.NewStyle().Bold(true)
	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Bold(true)
	okStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	missingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	warnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	keyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	leftStyle     = lipgloss.NewStyle().Padding(0, 1)
	rightStyle    = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, false, true)
	confirmStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 3)
)

type model struct {
	configs []Config
	cursor  int

	width, height int

	loading bool
	busy    bool
	status  string
	err     error

	confirm *Config // config dont la suppression attend confirmation
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
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case loadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.status = ""
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

	// Confirmation de suppression en attente.
	if m.confirm != nil {
		switch key {
		case "o", "y", "enter":
			return m.runRemove()
		case "n", "esc", "q", "ctrl+c":
			m.confirm = nil
		}
		return m, nil
	}

	if key == "q" || key == "ctrl+c" {
		return m, tea.Quit
	}
	if m.loading || m.busy {
		return m, nil
	}

	switch key {
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter", "a":
		return m.actionSelected()
	case "d":
		if c := m.current(); c != nil {
			if c.Installed {
				m.confirm = c
			} else {
				m.err = fmt.Errorf("%s : rien à supprimer (absent)", c.Branch)
			}
		}
	case "i":
		return m.globalAction(true)
	case "u":
		return m.globalAction(false)
	case "r":
		m.loading = true
		m.err = nil
		m.status = ""
		return m, discoverCmd()
	}
	return m, nil
}

// actionSelected installe ou met à jour la configuration sélectionnée.
func (m model) actionSelected() (tea.Model, tea.Cmd) {
	c := m.current()
	if c == nil {
		return m, nil
	}
	s, err := stepsFor(c)
	if err != nil {
		m.err = err
		return m, nil
	}
	return m.exec([]step{s}, nil)
}

// globalAction installe les manquantes (missing=true) ou met à jour toutes les
// installées (missing=false).
func (m model) globalAction(missing bool) (tea.Model, tea.Cmd) {
	var steps []step
	var warns []string
	for i := range m.configs {
		c := &m.configs[i]
		if c.Hidden {
			continue
		}
		if missing == c.Installed {
			continue
		}
		s, err := stepsFor(c)
		if err != nil {
			warns = append(warns, err.Error())
			continue
		}
		steps = append(steps, s)
	}
	return m.exec(steps, warns)
}

// runRemove exécute la suppression confirmée.
func (m model) runRemove() (tea.Model, tea.Cmd) {
	c := m.confirm
	m.confirm = nil
	if c == nil {
		return m, nil
	}
	s, err := removeStepFor(c)
	if err != nil {
		m.err = err
		return m, nil
	}
	return m.exec([]step{s}, nil)
}

// exec lance les étapes (avec d'éventuels avertissements) via tea.ExecProcess.
func (m model) exec(steps []step, warns []string) (tea.Model, tea.Cmd) {
	all := make([]step, 0, len(steps)+len(warns))
	for _, w := range warns {
		all = append(all, step{
			desc: "attention",
			sh:   "printf '\\033[33m%s\\033[0m\\n' " + shq("attention : "+w),
		})
	}
	all = append(all, steps...)

	if len(all) == 0 {
		m.status = "rien à faire"
		return m, nil
	}
	path, cleanup, err := buildScript(all)
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
	if m.width <= 0 {
		m.width = 100
	}
	if m.height <= 0 {
		m.height = 30
	}
	if m.confirm != nil {
		return m.confirmView()
	}

	header := titleStyle.Render("configs") + subtitleStyle.Render("  —  "+repoSlug)

	leftW := 34
	rightW := m.width - leftW - 6
	if rightW < 34 {
		rightW = 34
	}

	left := leftStyle.Width(leftW).Render(m.leftPane())
	right := rightStyle.Width(rightW).Render(m.rightPane())
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	content := lipgloss.JoinVertical(lipgloss.Left, header, "", body, "", m.footer())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) leftPane() string {
	vis := m.visible()
	if len(vis) == 0 {
		return dimStyle.Render("aucune configuration")
	}
	var b strings.Builder
	for i, ci := range vis {
		c := m.configs[ci]
		mark, _ := statusMark(&c)
		if i == m.cursor {
			fmt.Fprintf(&b, "%s %s\n", cursorStyle.Render("›"), mark+" "+selectedStyle.Render(c.label()))
		} else {
			fmt.Fprintf(&b, "  %s\n", mark+" "+c.label())
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) rightPane() string {
	c := m.current()
	if c == nil {
		return dimStyle.Render("Sélectionnez une configuration.")
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
	fmt.Fprintf(&b, "%s %s %s\n", titleStyle.Render("statut   "), word, mark)
	fmt.Fprintf(&b, "%s %s\n", titleStyle.Render("cible    "), c.TargetPath)
	fmt.Fprintf(&b, "%s %s\n", titleStyle.Render("gestion  "), managementWord(c))
	if c.Desync {
		b.WriteString(warnStyle.Render("désync : cible non gérée par git") + "\n")
	}
	if c.LastCommit != "" {
		fmt.Fprintf(&b, "%s %s\n", titleStyle.Render("commit   "), c.LastCommit)
	}
	if len(c.Prereq) > 0 {
		b.WriteString(titleStyle.Render("prérequis") + " ")
		for i, p := range c.Prereq {
			mk := okStyle.Render("✓")
			if !c.Prereqs[p] {
				mk = missingStyle.Render("✗")
			}
			if i > 0 {
				b.WriteString("  ")
			}
			b.WriteString(p + mk)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.actionsHint(c))
	return b.String()
}

// managementWord décrit comment la configuration est (ou serait) gérée, git ou
// install.sh — y compris quand elle n'est pas encore installée.
func managementWord(c *Config) string {
	switch {
	case c.GitManaged && c.GitBranch != "":
		return "dépôt git (branche " + c.GitBranch + ")"
	case c.GitManaged:
		return "dépôt git"
	case c.HasInstall:
		return "install.sh de la branche"
	case c.Installed:
		return "fichiers seuls"
	default:
		return "non installé"
	}
}

func (m model) actionsHint(c *Config) string {
	var parts []string
	if c.Installed {
		parts = append(parts, keyStyle.Render("a")+" mettre à jour")
		parts = append(parts, keyStyle.Render("d")+" supprimer")
	} else {
		parts = append(parts, keyStyle.Render("a")+" installer")
	}
	return dimStyle.Render("cette config :") + " " + strings.Join(parts, dimStyle.Render("  ·  "))
}

func (m model) footer() string {
	var b strings.Builder
	switch {
	case m.err != nil:
		b.WriteString(errStyle.Render(m.err.Error()) + "\n")
	case m.loading:
		b.WriteString(subtitleStyle.Render("actualisation…") + "\n")
	case m.status != "":
		b.WriteString(dimStyle.Render(m.status) + "\n")
	}
	b.WriteString(helpStyle.Render("↑/↓ sélectionner · entrée/a agir · d supprimer · i installer les manquants · u tout mettre à jour · r actualiser · q quitter"))
	return b.String()
}

func (m model) confirmView() string {
	c := m.confirm
	body := warnStyle.Render("Désinstaller "+c.label()+" ?") + "\n\n"
	if c.HasInstall {
		body += "Retrait via l'install.sh de la branche (--uninstall).\n"
	} else {
		body += "Suppression du dossier " + c.TargetPath + " (rm -rf).\n"
	}
	body += "\n" + keyStyle.Render("o") + " confirmer     " + keyStyle.Render("n") + " annuler"
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, confirmStyle.Render(body))
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
