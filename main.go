package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

const CONF_FILE = "tc.yaml"

type Target struct {
	Path    string   `yaml:"path"`
	Entries []string `yaml:"entries"`
}

type Copyjob struct {
	Targets []Target `yaml:"targets"`
}

var (
	copyjob Copyjob
	ctx     *kong.Context
)

type Context struct {
	Verbose    bool
	Buffersize string
}

type AddCmd struct {
	Target string   `required:"" short:"t" help:"target to copy to" type:"path"`
	Paths  []string `arg:"" optional:"" name:"path" help:"path to add to target" type:"path"`
}

type RmCmd struct {
	Target string   `required:"" short:"t" help:"target to remove from" type:"path"`
	Paths  []string `arg:"" optional:"" name:"path" help:"path to remove from target" type:"path"`
}

type RmTargetCmd struct {
	Target string `arg:"" help:"target to remove" type:"path"`
}

type LsCmd struct {
	Short bool `help:"short output" short:"s"`
}

type ClearCmd struct {
}

type RunCmd struct {
}

var cli struct {
	Verbose    bool        `help:"verbose output" short:"v"`
	Buffersize string      `help:"size of buffer to use for copy" short:"b" default:"64k"`
	Add        AddCmd      `cmd:"" help:"add files/folder to target"`
	Rm         RmCmd       `cmd:"" help:"remove files/folder from target"`
	RmTarget   RmTargetCmd `cmd:"rm-target" help:"remove target from copyjob"`
	Ls         LsCmd       `cmd:"" help:"list current copy job"`
	Clear      ClearCmd    `cmd:"" help:"empty current copy job"`
	Run        RunCmd      `cmd:"" help:"run copyjob"`
}

func (a *AddCmd) Run(c *Context) error {
	for i, t := range copyjob.Targets {
		if t.Path == a.Target {
			copyjob.Targets[i].Entries = append(copyjob.Targets[i].Entries, a.Paths...)
			if c.Verbose {
				fmt.Printf("added files to %s\n", t.Path)
			}
			return saveCopyjob()
		}
	}
	t := Target{
		Path:    a.Target,
		Entries: a.Paths,
	}
	copyjob.Targets = append(copyjob.Targets, t)
	if c.Verbose {
		fmt.Printf("added target %s\n", t.Path)
	}
	return saveCopyjob()
}

func (r *RmTargetCmd) Run(c *Context) error {
	for i, t := range copyjob.Targets {
		if t.Path == r.Target {
			copyjob.Targets = append(copyjob.Targets[:i], copyjob.Targets[i+1:]...)
			return saveCopyjob()
		}
	}
	return nil
}

func (r *RmCmd) Run(c *Context) error {
	for i, t := range copyjob.Targets {
		if t.Path == r.Target {
			for j, e := range t.Entries {
				if e == r.Paths[0] {
					copyjob.Targets[i].Entries = append(t.Entries[:j], t.Entries[j+1:]...)
					break
				}
			}
			return saveCopyjob()
		}
	}
	return nil
}

func (l *LsCmd) Run(c *Context) error {
	for _, t := range copyjob.Targets {
		fmt.Println(t.Path)
		if !l.Short {
			for _, e := range t.Entries {
				fmt.Println("  ", e)
			}
		}
	}
	return nil
}

func (l *ClearCmd) Run(c *Context) error {
	copyjob.Targets = []Target{}
	return saveCopyjob()
}

func (r *RunCmd) Run(c *Context) error {
	for _, t := range copyjob.Targets {
		c.copyTarget(&t)
	}
	return nil
}

func main() {
	_, err := os.Stat(CONF_FILE)
	if err == nil {
		err = loadCopyjob()
		if err != nil {
			log.Fatalf("error parsing %s: %v", CONF_FILE, err)
		}
	}
	ctx = kong.Parse(&cli)
	err = ctx.Run(&Context{Verbose: cli.Verbose})
	ctx.FatalIfErrorf(err)
}

func loadCopyjob() error {
	b, err := os.ReadFile(CONF_FILE)
	if err != nil {
		log.Fatalf("error reading %s: %v", CONF_FILE, err)
	}
	copyjob = Copyjob{}
	return yaml.Unmarshal(b, &copyjob)
}

func saveCopyjob() error {
	b, err := yaml.Marshal(&copyjob)
	if err != nil {
		return fmt.Errorf("error marshall copyjob: %v", err)
	}
	return os.WriteFile(CONF_FILE, b, 0644)
}
