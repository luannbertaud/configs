package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

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
	updateStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	missingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	warnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	keyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	leftStyle     = lipgloss.NewStyle().Padding(0, 1)
	rightStyle    = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, false, true)
	boxStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2)
	confirmStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 3)
)

// helpLine tient sur une seule ligne dans le cadre.
const helpLine = "↑/↓ choisir · entrée/a agir · d supprimer · i manquants · u màj · r actualiser · q quitter"

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type model struct {
	pkgs   []Package
	cursor int

	width, height int

	loading   bool
	busy      bool
	status    string
	opErr     string // erreur d'opération, conservée jusqu'au prochain refresh/action
	spinner   int
	spinnerID int

	confirm *Package
}

// messages.
type loadedMsg struct {
	pkgs []Package
	err  error
}

type finishedMsg struct{ err error }

type tickMsg struct{ id int }

func newModel() model {
	return model{loading: true, spinnerID: 1}
}

func runTUI() error {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func discoverCmd() tea.Cmd {
	return func() tea.Msg {
		pkgs, err := discover()
		return loadedMsg{pkgs, err}
	}
}

func tickCmd(id int) tea.Cmd {
	return tea.Tick(110*time.Millisecond, func(time.Time) tea.Msg { return tickMsg{id} })
}

// startLoading relance la découverte avec un spinner identifié, pour ignorer
// les ticks d'une actualisation précédente.
func (m *model) startLoading() tea.Cmd {
	m.loading = true
	m.spinnerID++
	return tea.Batch(discoverCmd(), tickCmd(m.spinnerID))
}

func (m model) Init() tea.Cmd {
	return tea.Batch(discoverCmd(), tickCmd(m.spinnerID))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		if !m.loading || msg.id != m.spinnerID {
			return m, nil
		}
		m.spinner++
		return m, tickCmd(msg.id)

	case loadedMsg:
		m.loading = false
		if msg.err != nil {
			m.opErr = msg.err.Error()
		} else {
			m.pkgs = msg.pkgs
			m.clampCursor()
		}
		return m, nil

	case finishedMsg:
		m.busy = false
		if msg.err != nil {
			m.opErr = "échec : " + msg.err.Error()
		} else {
			m.opErr = ""
			m.status = ""
		}
		return m, m.startLoading()

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
			return m.runUninstall()
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
		if p := m.current(); p != nil {
			if p.Installed {
				m.confirm = p
			} else {
				m.opErr = p.Branch + " : rien à supprimer (absent)"
			}
		}
	case "i":
		return m.globalAction(true)
	case "u":
		return m.globalAction(false)
	case "r":
		m.opErr = ""
		m.status = ""
		return m, m.startLoading()
	}
	return m, nil
}

// actionSelected installe, met à jour ou signale l'état du paquet sélectionné.
func (m model) actionSelected() (tea.Model, tea.Cmd) {
	p := m.current()
	if p == nil {
		return m, nil
	}
	var steps []step
	var err error
	switch p.State {
	case StateAbsent:
		steps, err = installSteps(p)
	case StateUpdate:
		steps, err = updateSteps(p)
	case StateUpToDate:
		m.status = p.Branch + " : déjà à jour"
		m.opErr = ""
		return m, nil
	default:
		err = fmt.Errorf("%s : paquet désync, intervention manuelle requise", p.Branch)
	}
	if err != nil {
		m.opErr = err.Error()
		return m, nil
	}
	return m.exec(steps, nil)
}

// globalAction installe les manquants (missing=true) ou met à jour tous les
// paquets à update (missing=false).
func (m model) globalAction(missing bool) (tea.Model, tea.Cmd) {
	var steps []step
	var warns []string
	for i := range m.pkgs {
		p := &m.pkgs[i]
		if p.Hidden {
			continue
		}
		if missing && p.State != StateAbsent {
			continue
		}
		if !missing && p.State != StateUpdate {
			continue
		}
		var s []step
		var err error
		if missing {
			s, err = installSteps(p)
		} else {
			s, err = updateSteps(p)
		}
		if err != nil {
			warns = append(warns, err.Error())
			continue
		}
		steps = append(steps, s...)
	}
	return m.exec(steps, warns)
}

// runUninstall exécute la suppression confirmée.
func (m model) runUninstall() (tea.Model, tea.Cmd) {
	p := m.confirm
	m.confirm = nil
	if p == nil {
		return m, nil
	}
	steps, err := uninstallSteps(p)
	if err != nil {
		m.opErr = err.Error()
		return m, nil
	}
	return m.exec(steps, nil)
}

// exec lance les étapes (avec d'éventuels avertissements) via tea.ExecProcess.
func (m model) exec(steps []step, warns []string) (tea.Model, tea.Cmd) {
	m.opErr = ""
	m.status = ""

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
		m.opErr = err.Error()
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

// visible renvoie les index des paquets affichables (hors HIDDEN).
func (m model) visible() []int {
	var idx []int
	for i := range m.pkgs {
		if !m.pkgs[i].Hidden {
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

func (m model) current() *Package {
	vis := m.visible()
	if m.cursor < 0 || m.cursor >= len(vis) {
		return nil
	}
	return &m.pkgs[vis[m.cursor]]
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

	contentW := m.contentWidth()
	if m.loading {
		return m.loaderView(contentW)
	}

	header := titleStyle.Render("configs") + subtitleStyle.Render("  —  "+repoSlug)

	leftW := 22
	rightW := contentW - leftW - 5
	if rightW < 24 {
		rightW = 24
	}

	left := leftStyle.Width(leftW).Render(m.leftPane())
	right := rightStyle.Width(rightW).Render(m.rightPane())
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	content := lipgloss.JoinVertical(lipgloss.Left, header, "", body, "", m.footer())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxStyle.Width(contentW+4).Render(content))
}

// contentWidth dimensionne le cadre pour que la ligne d'aide tienne sur une
// seule ligne, tout en gardant des marges sur les côtés.
func (m model) contentWidth() int {
	w := lipgloss.Width(helpLine)
	if max := m.width - 8; w > max {
		w = max
	}
	if w < 50 {
		w = 50
	}
	return w
}

func (m model) loaderView(contentW int) string {
	frame := spinnerFrames[m.spinner%len(spinnerFrames)]
	body := titleStyle.Render("configs") + "\n\n" + subtitleStyle.Render(frame+"  actualisation…")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		boxStyle.Width(contentW+4).Render(body))
}

func (m model) leftPane() string {
	vis := m.visible()
	var b strings.Builder
	b.WriteString(dimStyle.Render("PAQUETS") + "\n\n")
	if len(vis) == 0 {
		b.WriteString(dimStyle.Render("aucun paquet"))
		return b.String()
	}
	for i, pi := range vis {
		p := m.pkgs[pi]
		mark, _ := stateMark(p.State)
		if i == m.cursor {
			fmt.Fprintf(&b, "%s %s\n", cursorStyle.Render("›"), mark+" "+selectedStyle.Render(p.label()))
		} else {
			fmt.Fprintf(&b, "  %s\n", mark+" "+p.label())
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) rightPane() string {
	p := m.current()
	if p == nil {
		return dimStyle.Render("Sélectionnez un paquet.")
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render(p.label()))
	if p.Hidden {
		b.WriteString(dimStyle.Render(" (masqué)"))
	}
	b.WriteString("\n\n")
	if p.Description != "" {
		b.WriteString(p.Description + "\n")
	}
	b.WriteString("\n")

	mark, word := stateMark(p.State)
	fmt.Fprintf(&b, "%s %s %s\n", titleStyle.Render("état     "), word, mark)
	fmt.Fprintf(&b, "%s %s\n", titleStyle.Render("cible    "), p.TargetPath)
	fmt.Fprintf(&b, "%s %s\n", titleStyle.Render("gestion  "), managementWord(p))
	if p.LastCommit != "" {
		fmt.Fprintf(&b, "%s %s\n", titleStyle.Render("commit   "), p.LastCommit)
	}
	if len(p.Prereq) > 0 {
		b.WriteString(titleStyle.Render("prérequis") + " ")
		for i, c := range p.Prereq {
			mk := okStyle.Render("✓")
			if !p.Prereqs[c] {
				mk = missingStyle.Render("✗")
			}
			if i > 0 {
				b.WriteString("  ")
			}
			b.WriteString(c + mk)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.actionsHint(p))
	return b.String()
}

// managementWord décrit comment le paquet est géré localement.
func managementWord(p *Package) string {
	switch {
	case p.Managed:
		return "clone git · branche " + p.Branch
	case p.Installed && p.DesyncReason != "":
		return "désync · " + p.DesyncReason
	case p.Installed:
		return "présent, non géré"
	default:
		return "non installé"
	}
}

func (m model) actionsHint(p *Package) string {
	var parts []string
	switch p.State {
	case StateAbsent:
		parts = append(parts, keyStyle.Render("a")+" installer")
	case StateUpdate:
		parts = append(parts, keyStyle.Render("a")+" mettre à jour")
	case StateUpToDate:
		parts = append(parts, dimStyle.Render("déjà à jour"))
	default:
		parts = append(parts, dimStyle.Render("désync — intervention manuelle"))
	}
	if p.Installed {
		parts = append(parts, keyStyle.Render("d")+" supprimer")
	}
	return dimStyle.Render("cette config :") + " " + strings.Join(parts, dimStyle.Render("  ·  "))
}

func (m model) footer() string {
	var lines []string
	if m.opErr != "" {
		lines = append(lines, errStyle.Render(m.opErr))
	} else if m.status != "" {
		lines = append(lines, dimStyle.Render(m.status))
	}
	lines = append(lines, helpStyle.Render(helpLine))
	return strings.Join(lines, "\n")
}

func (m model) confirmView() string {
	p := m.confirm
	body := warnStyle.Render("Désinstaller "+p.label()+" ?") + "\n\n"
	if p.PostUninstall != "" && p.Managed {
		body += "Le script post_uninstall sera exécuté, puis le dossier sera retiré.\n"
	} else if p.Managed {
		body += "Le dossier " + p.TargetPath + " sera supprimé (rm -rf).\n"
	} else {
		body += "Dossier non géré : suppression brute de " + p.TargetPath + " (rm -rf),\nsans exécuter de script.\n"
	}
	body += "\n" + keyStyle.Render("o") + " confirmer     " + keyStyle.Render("n") + " annuler"
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, confirmStyle.Render(body))
}

// stateMark renvoie l'icône et le mot d'état d'un paquet.
func stateMark(s State) (string, string) {
	switch s {
	case StateAbsent:
		return missingStyle.Render("✗"), "absent"
	case StateDesync:
		return warnStyle.Render("!"), "désync"
	case StateUpdate:
		return updateStyle.Render("↑"), "à update"
	default:
		return okStyle.Render("✓"), "à jour"
	}
}
