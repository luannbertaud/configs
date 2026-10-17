package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Emplacement du dépôt et des fichiers servis en brut.
const (
	repoSlug     = "luannbertaud/configs"
	repoURL      = "https://github.com/luannbertaud/configs.git"
	rawBase      = "https://raw.githubusercontent.com/luannbertaud/configs"
	manifestFile = "configs.json"
	installFile  = "install.sh"
)

// Manifest est le descripteur d'une branche config, lu depuis configs.json à la
// racine de la branche.
type Manifest struct {
	Name   string   `json:"name,omitempty"`
	Desc   string   `json:"desc"`
	Target string   `json:"target"`
	Prereq []string `json:"prereq,omitempty"`
	Hidden bool     `json:"hidden,omitempty"`
}

// Config est un manifeste enrichi de l'état de la machine courante.
type Config struct {
	Branch string
	Manifest

	TargetPath string
	Installed  bool            // le dossier cible existe
	GitManaged bool            // clone git dont l'origin est ce dépôt
	GitBranch  string          // branche actuellement extraite dans la cible
	Desync     bool            // cible présente mais pas un clone de la branche
	HasInstall bool            // la branche embarque un install.sh
	Prereqs    map[string]bool // commande -> présente ?
	LastCommit string
}

// label renvoie le nom affiché (name du manifeste, sinon la branche).
func (c *Config) label() string {
	if c.Name != "" {
		return c.Name
	}
	return c.Branch
}

// listBranches énumère les branches distantes via `git ls-remote --heads`.
func listBranches() ([]string, error) {
	out, err := exec.Command("git", "ls-remote", "--heads", repoURL).Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-remote: %w", err)
	}
	var branches []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		if ref, ok := strings.CutPrefix(fields[1], "refs/heads/"); ok {
			branches = append(branches, ref)
		}
	}
	return branches, nil
}

// fetchManifest télécharge le configs.json d'une branche. Il renvoie (nil, nil)
// quand la branche n'en a pas (ce n'est donc pas une branche config).
func fetchManifest(client *http.Client, branch string) (*Manifest, error) {
	url := fmt.Sprintf("%s/%s/%s", rawBase, branch, manifestFile)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		var m Manifest
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			return nil, fmt.Errorf("%s: %w", url, err)
		}
		return &m, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, fmt.Errorf("%s: HTTP %d", url, resp.StatusCode)
	}
}

// remoteExists teste l'existence d'un fichier servi en brut (HEAD).
func remoteExists(client *http.Client, url string) bool {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// discover liste les branches, ne garde que celles portant un configs.json et
// renvoie leur état sur cette machine, trié par nom de branche.
func discover() ([]Config, error) {
	client := &http.Client{Timeout: 20 * time.Second}

	branches, err := listBranches()
	if err != nil {
		return nil, err
	}

	type result struct {
		cfg *Config
		err error
	}
	results := make([]result, len(branches))
	var wg sync.WaitGroup
	for i, branch := range branches {
		wg.Add(1)
		go func(i int, branch string) {
			defer wg.Done()
			m, err := fetchManifest(client, branch)
			if err != nil {
				results[i] = result{nil, err}
				return
			}
			if m == nil {
				return // pas une branche config
			}
			c := &Config{Branch: branch, Manifest: *m}
			c.inspect(client)
			results[i] = result{c, nil}
		}(i, branch)
	}
	wg.Wait()

	var configs []Config
	var errs []string
	for _, r := range results {
		if r.err != nil {
			errs = append(errs, r.err.Error())
			continue
		}
		if r.cfg != nil {
			configs = append(configs, *r.cfg)
		}
	}
	sort.Slice(configs, func(i, j int) bool { return configs[i].Branch < configs[j].Branch })

	if len(configs) == 0 && len(errs) > 0 {
		return nil, fmt.Errorf("découverte impossible: %s", strings.Join(errs, "; "))
	}
	return configs, nil
}

// inspect renseigne l'état local (cible, git, install.sh, pré-requis).
func (c *Config) inspect(client *http.Client) {
	c.TargetPath = expandTarget(c.Target)
	c.Prereqs = map[string]bool{}
	for _, p := range c.Prereq {
		_, err := exec.LookPath(p)
		c.Prereqs[p] = err == nil
	}

	if fi, err := os.Stat(c.TargetPath); err == nil && fi.IsDir() {
		c.Installed = true
		if isGitWorktree(c.TargetPath) && originIsRepo(c.TargetPath) {
			c.GitManaged = true
			if out, err := gitOut(c.TargetPath, "branch", "--show-current"); err == nil {
				c.GitBranch = strings.TrimSpace(out)
			}
			if out, err := gitOut(c.TargetPath, "log", "-1", "--format=%h %s"); err == nil {
				c.LastCommit = strings.TrimSpace(out)
			}
		} else {
			c.Desync = true
		}
	}

	c.HasInstall = remoteExists(client, fmt.Sprintf("%s/%s/%s", rawBase, c.Branch, installFile))
}

func isGitWorktree(dir string) bool {
	out, err := gitOut(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

func originIsRepo(dir string) bool {
	out, err := gitOut(dir, "remote", "get-url", "origin")
	if err != nil {
		return false
	}
	u := strings.TrimSpace(out)
	return u == repoURL || strings.Contains(u, repoSlug)
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

// stepsFor construit l'étape d'installation ou de mise à jour d'une config :
// délégation au install.sh de la branche s'il existe, sinon git.
func stepsFor(c *Config) (step, error) {
	switch {
	case c.HasInstall:
		label := "Installer " + c.Branch
		if c.Installed {
			label = "Mettre à jour " + c.Branch
		}
		url := fmt.Sprintf("%s/%s/%s", rawBase, c.Branch, installFile)
		return step{desc: label + " (install.sh)", sh: "curl -fsSL " + shq(url) + " | sh"}, nil

	case !c.Installed:
		if err := os.MkdirAll(filepath.Dir(c.TargetPath), 0o755); err != nil {
			return step{}, err
		}
		return step{
			desc: "Installer " + c.Branch,
			sh: fmt.Sprintf("git clone --branch %s --single-branch %s %s",
				shq(c.Branch), shq(repoURL), shq(c.TargetPath)),
		}, nil

	case c.GitManaged:
		return step{
			desc: "Mettre à jour " + c.Branch,
			sh:   fmt.Sprintf("git -C %s pull --ff-only", shq(c.TargetPath)),
		}, nil

	default:
		return step{}, fmt.Errorf("%s : %s existe mais n'est pas géré par git (désync)", c.Branch, c.TargetPath)
	}
}

// removeStepFor construit l'étape de suppression d'une configuration :
// délégation à `install.sh --uninstall` s'il existe, sinon suppression du
// dossier cible quand c'est un clone git géré par ce dépôt.
func removeStepFor(c *Config) (step, error) {
	if !c.Installed {
		return step{}, fmt.Errorf("%s : rien à supprimer (absent)", c.Branch)
	}
	if c.HasInstall {
		url := fmt.Sprintf("%s/%s/%s", rawBase, c.Branch, installFile)
		return step{
			desc: "Désinstaller " + c.Branch + " (install.sh --uninstall)",
			sh:   "curl -fsSL " + shq(url) + " | sh -s -- --uninstall",
		}, nil
	}
	if c.GitManaged {
		return step{
			desc: "Désinstaller " + c.Branch,
			sh:   "rm -rf -- " + shq(c.TargetPath),
		}, nil
	}
	return step{}, fmt.Errorf("%s : %s non géré par git (désync) — suppression manuelle", c.Branch, c.TargetPath)
}

// buildScript écrit un script shell temporaire exécutant les étapes dans
// l'ordre et renvoie son chemin ainsi qu'une fonction de nettoyage.
func buildScript(steps []step) (string, func(), error) {
	dir, err := os.MkdirTemp("", "configs-*")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(dir) }

	var b strings.Builder
	b.WriteString("#!/bin/sh\nset -eu\n")
	for _, s := range steps {
		fmt.Fprintf(&b, "\nprintf '\\n=== %%s ===\\n' %s\n", shq(s.desc))
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

// findConfig retrouve une config par nom de branche.
func findConfig(configs []Config, name string) *Config {
	for i := range configs {
		if configs[i].Branch == name || configs[i].Name == name {
			return &configs[i]
		}
	}
	return nil
}
