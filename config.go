package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Emplacement du dépôt (codé en dur) et du manifeste porté par chaque branche.
// repoSlug/repoURL sont des variables pour pouvoir être remplacées au build de
// test (go build -ldflags "-X main.repoURL=...").
var (
	repoSlug = "luannbertaud/configs"
	repoURL  = "https://github.com/luannbertaud/configs.git"
)

const manifestFile = "manifest.json"

// Manifest est le descripteur d'un paquet, lu depuis manifest.json à la racine
// de sa branche.
type Manifest struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description"`
	Target        string   `json:"target"`
	Prereq        []string `json:"prereq,omitempty"`
	PostInstall   string   `json:"post_install,omitempty"`
	PostUninstall string   `json:"post_uninstall,omitempty"`
	Hidden        bool     `json:"hidden,omitempty"`
}

// State décrit l'état d'un paquet sur la machine.
type State string

const (
	StateAbsent   State = "absent"
	StateDesync   State = "desync"
	StateUpdate   State = "update"
	StateUpToDate State = "uptodate"
)

// Package est un manifeste enrichi de l'état local.
type Package struct {
	Branch string
	Manifest

	TargetPath string
	RemoteTip  string

	Installed  bool
	Managed    bool
	Head       string
	LastCommit string
	Prereqs    map[string]bool

	State State
}

// label renvoie le nom affiché (name du manifeste, sinon la branche).
func (p *Package) label() string {
	if p.Name != "" {
		return p.Name
	}
	return p.Branch
}

// ----- découverte -----------------------------------------------------------

// discover récupère un clone partiel éphémère de toutes les branches (un seul
// commit par branche, sans blobs), lit les manifestes à la demande puis calcule
// l'état de chaque paquet sur cette machine.
func discover() ([]Package, error) {
	dir, cleanup, err := cloneMeta()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	tips, err := branchTips(dir)
	if err != nil {
		return nil, err
	}

	var pkgs []Package
	for branch, tip := range tips {
		// Vérifie l'existence du manifeste dans l'arbre (déjà présent avec
		// --filter=blob:none) sans récupérer son contenu.
		has, _ := gitOut(dir, "ls-tree", "--name-only", "origin/"+branch, "--", manifestFile)
		if strings.TrimSpace(has) == "" {
			continue
		}
		data, err := gitOut(dir, "cat-file", "-p", "origin/"+branch+":"+manifestFile)
		if err != nil {
			continue
		}
		var m Manifest
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			continue // manifeste illisible : on ignore la branche
		}
		if strings.TrimSpace(m.Target) == "" {
			continue // manifeste invalide : cible obligatoire
		}
		p := Package{Branch: branch, Manifest: m, RemoteTip: tip}
		p.inspect()
		pkgs = append(pkgs, p)
	}
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].Branch < pkgs[j].Branch })
	return pkgs, nil
}

// cloneMeta crée un clone partiel et éphémère de toutes les branches :
// --depth 1 (un commit), --no-single-branch (toutes les branches),
// --filter=blob:none (contenu des fichiers à la demande), --no-checkout.
// Repli sans filtre si le serveur ne le supporte pas.
func cloneMeta() (string, func(), error) {
	dir, err := os.MkdirTemp("", "configs-meta-*")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(dir) }

	clone := func(extra ...string) error {
		args := []string{"clone", "--quiet", "--depth", "1", "--no-single-branch", "--no-checkout"}
		args = append(args, extra...)
		args = append(args, repoURL, dir)
		cmd := exec.Command("git", args...)
		cmd.Stderr = nil
		return cmd.Run()
	}

	if err := clone("--filter=blob:none"); err != nil {
		// Certains serveurs refusent le filtre partiel : on retente sans.
		os.RemoveAll(dir)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", nil, err
		}
		if err := clone(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("découverte des paquets: %w", err)
		}
	}
	return dir, cleanup, nil
}

// branchTips renvoie, pour chaque branche distante, le SHA de son tip.
func branchTips(dir string) (map[string]string, error) {
	out, err := gitOut(dir, "for-each-ref", "--format=%(refname:short) %(objectname)", "refs/remotes/origin")
	if err != nil {
		return nil, fmt.Errorf("liste des branches: %w", err)
	}
	tips := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		branch, ok := strings.CutPrefix(fields[0], "origin/")
		if !ok || branch == "HEAD" || strings.Contains(branch, "->") {
			continue
		}
		tips[branch] = fields[1]
	}
	return tips, nil
}

// inspect renseigne l'état local du paquet.
func (p *Package) inspect() {
	p.TargetPath = expandTarget(p.Target)
	p.Prereqs = map[string]bool{}
	for _, c := range p.Prereq {
		_, err := exec.LookPath(c)
		p.Prereqs[c] = err == nil
	}

	fi, err := os.Stat(p.TargetPath)
	if err != nil || !fi.IsDir() {
		p.State = StateAbsent
		return
	}
	p.Installed = true

	if !isGitWorktree(p.TargetPath) || !originMatches(p.TargetPath) || currentBranch(p.TargetPath) != p.Branch {
		p.State = StateDesync
		return
	}

	p.Managed = true
	if h, err := gitOut(p.TargetPath, "rev-parse", "HEAD"); err == nil {
		p.Head = strings.TrimSpace(h)
	}
	if out, err := gitOut(p.TargetPath, "log", "-1", "--format=%h %s"); err == nil {
		p.LastCommit = strings.TrimSpace(out)
	}
	if isDirty(p.TargetPath) {
		p.State = StateDesync
		return
	}
	if p.Head != "" && p.Head != p.RemoteTip {
		p.State = StateUpdate
		return
	}
	p.State = StateUpToDate
}

func isGitWorktree(dir string) bool {
	out, err := gitOut(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// originMatches accepte les URLs SSH et HTTPS du dépôt (l'origin n'est jamais
// utilisé pour les opérations, seulement pour classifier le dossier).
func originMatches(dir string) bool {
	out, err := gitOut(dir, "remote", "get-url", "origin")
	if err != nil {
		return false
	}
	return strings.Contains(strings.TrimSpace(out), repoSlug)
}

func currentBranch(dir string) string {
	out, err := gitOut(dir, "branch", "--show-current")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// isDirty ne regarde que les fichiers suivis (les fichiers générés non suivis
// sont ignorés).
func isDirty(dir string) bool {
	out, err := gitOut(dir, "status", "--porcelain", "--untracked-files=no")
	return err == nil && strings.TrimSpace(out) != ""
}

func gitOut(dir string, args ...string) (string, error) {
	full := append([]string{"-C", dir}, args...)
	out, err := exec.Command("git", full...).Output()
	return string(out), err
}

// expandTarget résout la cible : "~" et "~/..." par rapport au HOME, les chemins
// absolus tels quels, sinon relatifs au HOME.
func expandTarget(s string) string {
	if s == "" {
		return ""
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return s
	}
	switch {
	case s == "~":
		return home
	case strings.HasPrefix(s, "~/"):
		return filepath.Join(home, s[2:])
	case filepath.IsAbs(s):
		return s
	default:
		return filepath.Join(home, s)
	}
}

// ----- plan d'exécution ------------------------------------------------------

// step est une commande shell accompagnée d'un libellé.
type step struct {
	desc string
	sh   string
}

// installSteps installe un paquet absent (clone mono-branche de la branche
// complète) puis exécute son post_install.
func installSteps(p *Package) ([]step, error) {
	if p.State != StateAbsent {
		return updateSteps(p)
	}
	if err := os.MkdirAll(filepath.Dir(p.TargetPath), 0o755); err != nil {
		return nil, err
	}
	steps := []step{{
		desc: "Installer " + p.Branch,
		sh: fmt.Sprintf("git clone --single-branch --branch %s %s %s",
			shq(p.Branch), shq(repoURL), shq(p.TargetPath)),
	}}
	if p.PostInstall != "" {
		steps = append(steps, postInstallStep(p, ""))
	}
	return steps, nil
}

// updateSteps met à jour un paquet (pull --ff-only via l'URL canonique) puis
// exécute son post_install si HEAD a changé.
func updateSteps(p *Package) ([]step, error) {
	switch p.State {
	case StateAbsent:
		return installSteps(p)
	case StateUpToDate:
		return nil, nil
	case StateUpdate:
		steps := []step{{
			desc: "Mettre à jour " + p.Branch,
			sh: fmt.Sprintf("( cd %s && git pull --ff-only %s %s )",
				shq(p.TargetPath), shq(repoURL), shq(p.Branch)),
		}}
		if p.PostInstall != "" {
			steps = append(steps, postInstallStep(p, p.Head))
		}
		return steps, nil
	default:
		return nil, fmt.Errorf("%s : état %s, mise à jour impossible", p.Branch, p.State)
	}
}

// uninstallSteps exécute post_uninstall (s'il existe dans la cible) puis retire
// le dossier cible.
func uninstallSteps(p *Package) ([]step, error) {
	if !p.Installed {
		return nil, fmt.Errorf("%s : rien à supprimer (absent)", p.Branch)
	}
	var steps []step
	if p.PostUninstall != "" {
		steps = append(steps, step{
			desc: "post-uninstall " + p.Branch,
			sh: fmt.Sprintf("if [ -f %s ]; then ( cd %s && CONFIGS_TARGET=%s sh %s ); fi",
				shq(filepath.Join(p.TargetPath, p.PostUninstall)),
				shq(p.TargetPath), shq(p.TargetPath), shq(p.PostUninstall)),
		})
	}
	steps = append(steps, step{
		desc: "Désinstaller " + p.Branch,
		sh:   "rm -rf -- " + shq(p.TargetPath),
	})
	return steps, nil
}

// postInstallStep exécute post_install avec cwd = cible. Si prev est fourni
// (update), l'étape n'est exécutée que si HEAD a changé. CONFIGS_COMMIT est lu
// après l'opération pour rester exact.
func postInstallStep(p *Package, prev string) step {
	env := fmt.Sprintf(`CONFIGS_TARGET=%s CONFIGS_COMMIT="$(git -C %s rev-parse HEAD)"`,
		shq(p.TargetPath), shq(p.TargetPath))
	if prev != "" {
		env += " CONFIGS_PREV=" + shq(prev)
	}
	body := fmt.Sprintf("( cd %s && %s sh %s )", shq(p.TargetPath), env, shq(p.PostInstall))
	if prev != "" {
		body = fmt.Sprintf("if [ %s != \"$(git -C %s rev-parse HEAD)\" ]; then %s; fi",
			shq(prev), shq(p.TargetPath), body)
	}
	return step{
		desc: "post-install " + p.Branch,
		sh:   body,
	}
}

// buildScript écrit un script shell temporaire exécutant les étapes dans
// l'ordre, en journalisant la commande exacte, et renvoie son chemin ainsi
// qu'une fonction de nettoyage.
func buildScript(steps []step) (string, func(), error) {
	dir, err := os.MkdirTemp("", "configs-*")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(dir) }

	var b strings.Builder
	b.WriteString("#!/bin/sh\nset -eu\n")
	for _, s := range steps {
		fmt.Fprintf(&b, "\nprintf '\\n\\033[1;36m=== %%s ===\\033[0m\\n' %s\n", shq(s.desc))
		fmt.Fprintf(&b, "printf '\\033[2m$ %%s\\033[0m\\n' %s\n", shq(s.sh))
		b.WriteString(s.sh)
		b.WriteString("\n")
	}

	path := filepath.Join(dir, "run.sh")
	if err := os.WriteFile(path, []byte(b.String()), 0o700); err != nil {
		cleanup()
		return "", nil, err
	}
	return path, cleanup, nil
}

// shq met une chaîne entre guillemets simples pour sh.
func shq(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// findPackage retrouve un paquet par nom de branche ou nom affiché.
func findPackage(pkgs []Package, name string) *Package {
	for i := range pkgs {
		if pkgs[i].Branch == name || pkgs[i].Name == name {
			return &pkgs[i]
		}
	}
	return nil
}
