// Command configs gère, depuis une seule branche main, les configurations
// versionnées du dépôt luannbertaud/configs : il liste celles qui sont
// installées sur la machine, les met à jour et installe les manquantes.
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
		configs, err := discover()
		if err != nil {
			fatal(err)
		}
		if opts.json {
			printJSON(configs)
		} else {
			printList(configs)
		}
	case opts.install != "":
		configs, err := discover()
		if err != nil {
			fatal(err)
		}
		c := findConfig(configs, opts.install)
		if c == nil {
			fatal(fmt.Errorf("configuration inconnue: %s", opts.install))
		}
		s, err := stepsFor(c)
		if err != nil {
			fatal(err)
		}
		if err := runSteps([]step{s}); err != nil {
			fatal(err)
		}
	case opts.uninstall != "":
		configs, err := discover()
		if err != nil {
			fatal(err)
		}
		c := findConfig(configs, opts.uninstall)
		if c == nil {
			fatal(fmt.Errorf("configuration inconnue: %s", opts.uninstall))
		}
		s, err := removeStepFor(c)
		if err != nil {
			fatal(err)
		}
		if err := runSteps([]step{s}); err != nil {
			fatal(err)
		}
	case opts.updateSet:
		configs, err := discover()
		if err != nil {
			fatal(err)
		}
		if opts.updateName != "" {
			c := findConfig(configs, opts.updateName)
			if c == nil {
				fatal(fmt.Errorf("configuration inconnue: %s", opts.updateName))
			}
			s, err := stepsFor(c)
			if err != nil {
				fatal(err)
			}
			if err := runSteps([]step{s}); err != nil {
				fatal(err)
			}
			return
		}
		steps, warns := collectUpdate(configs)
		if err := runSteps(steps, warns...); err != nil {
			fatal(err)
		}
	case opts.installMissing:
		configs, err := discover()
		if err != nil {
			fatal(err)
		}
		steps, warns := collectMissing(configs)
		if err := runSteps(steps, warns...); err != nil {
			fatal(err)
		}
	default:
		if !isTerminal() {
			// Sans TTY (pipe), on retombe sur la liste.
			configs, err := discover()
			if err != nil {
				fatal(err)
			}
			printList(configs)
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
				return o, errors.New("--install requiert un nom de configuration")
			}
			o.install = args[i]
		case len(a) > len("--install=") && a[:len("--install=")] == "--install=":
			o.install = a[len("--install="):]
		case a == "--uninstall":
			i++
			if i >= len(args) {
				return o, errors.New("--uninstall requiert un nom de configuration")
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
	fmt.Fprint(w, `configs — gère les configurations du dépôt luannbertaud/configs.

Usage :
  configs                      lance le TUI interactif
  configs --list [--json]      liste les configurations et leur état
  configs --install <nom>      installe une configuration (déjà installée : met à jour)
  configs --install-missing    installe toutes les configurations manquantes
  configs --update [nom]       met à jour une configuration, ou toutes si omis
  configs --uninstall <nom>    désinstalle une configuration
  configs --version            affiche la version

Le TUI propose : navigation ↑/↓, agir sur la sélection (entrée/a), supprimer (d),
installer les manquants (i), tout mettre à jour (u), actualiser (r), quitter (q).
`)
}

// ----- sorties --------------------------------------------------------------

func printList(configs []Config) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "NOM\tÉTAT\tCIBLE\tDESCRIPTION")
	for _, c := range configs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", c.label(), statusWord(&c), c.TargetPath, c.Desc)
	}
	tw.Flush()
	if len(configs) == 0 {
		fmt.Fprintln(os.Stderr, "configs: aucune configuration trouvée")
	}
}

func printJSON(configs []Config) {
	type row struct {
		Branch     string          `json:"branch"`
		Name       string          `json:"name"`
		Desc       string          `json:"desc"`
		Target     string          `json:"target"`
		Installed  bool            `json:"installed"`
		GitManaged bool            `json:"git_managed"`
		GitBranch  string          `json:"git_branch,omitempty"`
		Desync     bool            `json:"desync"`
		Hidden     bool            `json:"hidden"`
		Prereq     map[string]bool `json:"prereq,omitempty"`
	}
	out := make([]row, 0, len(configs))
	for _, c := range configs {
		out = append(out, row{
			Branch: c.Branch, Name: c.label(), Desc: c.Desc, Target: c.TargetPath,
			Installed: c.Installed, GitManaged: c.GitManaged, GitBranch: c.GitBranch,
			Desync: c.Desync, Hidden: c.Hidden, Prereq: c.Prereqs,
		})
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fatal(err)
	}
}

func statusWord(c *Config) string {
	switch {
	case !c.Installed:
		return "absent"
	case c.Desync:
		return "désync"
	default:
		return "installé"
	}
}

// ----- actions globales -----------------------------------------------------

// collectUpdate rassemble les étapes de mise à jour de toutes les configs
// installées. Les configs en désync sont signalées et ignorées.
func collectUpdate(configs []Config) ([]step, []string) {
	var steps []step
	var warns []string
	for i := range configs {
		c := &configs[i]
		if c.Hidden || !c.Installed {
			continue
		}
		s, err := stepsFor(c)
		if err != nil {
			warns = append(warns, err.Error())
			continue
		}
		steps = append(steps, s)
	}
	return steps, warns
}

// collectMissing rassemble les étapes d'installation des configs absentes.
func collectMissing(configs []Config) ([]step, []string) {
	var steps []step
	var warns []string
	for i := range configs {
		c := &configs[i]
		if c.Hidden || c.Installed {
			continue
		}
		s, err := stepsFor(c)
		if err != nil {
			warns = append(warns, err.Error())
			continue
		}
		steps = append(steps, s)
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
