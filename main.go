// Command configs gère les paquets du dépôt luannbertaud/configs : chaque
// branche portant un manifest.json est un paquet installé comme un clone git
// mono-branche et éventuellement finalisé par un script post_install.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"
)

// version est renseignée à la compilation (release).
var version = "dev"

func main() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "configs:", err)
		os.Exit(2)
	}
	if opts.help {
		printUsage(os.Stdout)
		return
	}
	if opts.version {
		fmt.Println("configs", version)
		return
	}

	switch {
	case opts.list:
		pkgs, err := discover()
		if err != nil {
			fatal(err)
		}
		if opts.json {
			printJSON(pkgs)
		} else {
			printList(pkgs)
		}
	case opts.install != "":
		pkgs, err := discover()
		if err != nil {
			fatal(err)
		}
		p := findPackage(pkgs, opts.install)
		if p == nil {
			fatal(fmt.Errorf("paquet inconnu: %s", opts.install))
		}
		steps, err := installSteps(p)
		if err != nil {
			fatal(err)
		}
		if err := runSteps(steps); err != nil {
			fatal(err)
		}
	case opts.installMissing:
		pkgs, err := discover()
		if err != nil {
			fatal(err)
		}
		steps, warns := collectMissing(pkgs)
		if err := runSteps(steps, warns...); err != nil {
			fatal(err)
		}
	case opts.updateSet:
		pkgs, err := discover()
		if err != nil {
			fatal(err)
		}
		if opts.updateName != "" {
			p := findPackage(pkgs, opts.updateName)
			if p == nil {
				fatal(fmt.Errorf("paquet inconnu: %s", opts.updateName))
			}
			steps, err := updateSteps(p)
			if err != nil {
				fatal(err)
			}
			if err := runSteps(steps); err != nil {
				fatal(err)
			}
			return
		}
		steps, warns := collectUpdate(pkgs)
		if err := runSteps(steps, warns...); err != nil {
			fatal(err)
		}
	case opts.uninstall != "":
		pkgs, err := discover()
		if err != nil {
			fatal(err)
		}
		p := findPackage(pkgs, opts.uninstall)
		if p == nil {
			fatal(fmt.Errorf("paquet inconnu: %s", opts.uninstall))
		}
		if !p.Managed {
			fmt.Fprintf(os.Stderr, "configs: attention: %s n'est pas géré, suppression brute sans script\n", p.Branch)
		}
		steps, err := uninstallSteps(p)
		if err != nil {
			fatal(err)
		}
		if err := runSteps(steps); err != nil {
			fatal(err)
		}
	default:
		if !isTerminal() {
			pkgs, err := discover()
			if err != nil {
				fatal(err)
			}
			printList(pkgs)
			return
		}
		if err := runTUI(); err != nil {
			fatal(err)
		}
	}
}

// ----- options --------------------------------------------------------------

type options struct {
	list           bool
	json           bool
	install        string
	uninstall      string
	updateSet      bool
	updateName     string
	installMissing bool
	help           bool
	version        bool
}

// parseArgs analyse argv. --update accepte un nom facultatif (ex. --update nvim).
func parseArgs(args []string) (options, error) {
	var o options
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--list":
			o.list = true
		case a == "--json":
			o.json = true
		case a == "--install-missing":
			o.installMissing = true
		case a == "--help" || a == "-h":
			o.help = true
		case a == "--version":
			o.version = true
		case a == "--install":
			i++
			if i >= len(args) {
				return o, errors.New("--install requiert un nom de paquet")
			}
			o.install = args[i]
		case len(a) > len("--install=") && a[:len("--install=")] == "--install=":
			o.install = a[len("--install="):]
		case a == "--uninstall":
			i++
			if i >= len(args) {
				return o, errors.New("--uninstall requiert un nom de paquet")
			}
			o.uninstall = args[i]
		case len(a) > len("--uninstall=") && a[:len("--uninstall=")] == "--uninstall=":
			o.uninstall = a[len("--uninstall="):]
		case a == "--update":
			o.updateSet = true
			if i+1 < len(args) && len(args[i+1]) > 0 && args[i+1][0] != '-' {
				i++
				o.updateName = args[i]
			}
		case len(a) > len("--update=") && a[:len("--update=")] == "--update=":
			o.updateSet = true
			o.updateName = a[len("--update="):]
		default:
			return o, fmt.Errorf("argument inconnu: %s", a)
		}
	}
	if o.install != "" && o.installMissing {
		return o, errors.New("--install et --install-missing sont exclusifs")
	}
	return o, nil
}

func printUsage(w *os.File) {
	fmt.Fprint(w, `configs — gère les paquets du dépôt luannbertaud/configs.

Chaque branche portant un manifest.json est un paquet. Les paquets sont
installés comme des clones git mono-branche (branche complète, sans les autres
branches).

Usage :
  configs                      lance le TUI interactif
  configs --list [--json]      liste les paquets et leur état
  configs --install <nom>      installe un paquet (ou le met à jour)
  configs --install-missing    installe tous les paquets absents
  configs --update [nom]       met à jour un paquet, ou tous ceux à jour à faire
  configs --uninstall <nom>    désinstalle un paquet
  configs --version            affiche la version

États : absent · désync · à update · à jour.
Le TUI propose : navigation ↑/↓, agir sur la sélection (entrée/a), supprimer (d),
installer les manquants (i), tout mettre à jour (u), actualiser (r), quitter (q).
`)
}

// ----- sorties --------------------------------------------------------------

func printList(pkgs []Package) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "NOM\tÉTAT\tCIBLE\tDESCRIPTION")
	for _, p := range pkgs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", p.label(), stateWord(p.State), p.TargetPath, p.Description)
	}
	tw.Flush()
	if len(pkgs) == 0 {
		fmt.Fprintln(os.Stderr, "configs: aucun paquet trouvé")
	}
}

func printJSON(pkgs []Package) {
	type row struct {
		Branch      string          `json:"branch"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Target      string          `json:"target"`
		State       string          `json:"state"`
		Installed   bool            `json:"installed"`
		Managed     bool            `json:"managed"`
		Head        string          `json:"head,omitempty"`
		RemoteTip   string          `json:"remote_tip,omitempty"`
		Hidden      bool            `json:"hidden"`
		PostInstall string          `json:"post_install,omitempty"`
		Prereq      map[string]bool `json:"prereq,omitempty"`
	}
	out := make([]row, 0, len(pkgs))
	for _, p := range pkgs {
		out = append(out, row{
			Branch: p.Branch, Name: p.label(), Description: p.Description,
			Target: p.TargetPath, State: string(p.State), Installed: p.Installed,
			Managed: p.Managed, Head: p.Head, RemoteTip: p.RemoteTip,
			Hidden: p.Hidden, PostInstall: p.PostInstall, Prereq: p.Prereqs,
		})
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fatal(err)
	}
}

func stateWord(s State) string {
	switch s {
	case StateAbsent:
		return "absent"
	case StateDesync:
		return "désync"
	case StateUpdate:
		return "à update"
	case StateUpToDate:
		return "à jour"
	}
	return string(s)
}

// ----- actions globales -----------------------------------------------------

// collectUpdate rassemble les étapes de mise à jour des paquets à update.
func collectUpdate(pkgs []Package) ([]step, []string) {
	var steps []step
	var warns []string
	for i := range pkgs {
		p := &pkgs[i]
		if p.Hidden || p.State != StateUpdate {
			continue
		}
		s, err := updateSteps(p)
		if err != nil {
			warns = append(warns, err.Error())
			continue
		}
		steps = append(steps, s...)
	}
	return steps, warns
}

// collectMissing rassemble les étapes d'installation des paquets absents.
func collectMissing(pkgs []Package) ([]step, []string) {
	var steps []step
	var warns []string
	for i := range pkgs {
		p := &pkgs[i]
		if p.Hidden || p.State != StateAbsent {
			continue
		}
		s, err := installSteps(p)
		if err != nil {
			warns = append(warns, err.Error())
			continue
		}
		steps = append(steps, s...)
	}
	return steps, warns
}

// runSteps exécute les étapes dans un script shell temporaire, en héritant des
// entrées/sorties standards.
func runSteps(steps []step, warns ...string) error {
	for _, w := range warns {
		fmt.Fprintln(os.Stderr, "configs: attention:", w)
	}
	if len(steps) == 0 {
		fmt.Fprintln(os.Stderr, "configs: rien à faire")
		return nil
	}
	path, cleanup, err := buildScript(steps)
	if err != nil {
		return err
	}
	defer cleanup()

	cmd := exec.Command("sh", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// isTerminal indique si la sortie standard est un terminal.
func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "configs:", err)
	os.Exit(1)
}
