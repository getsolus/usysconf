// Copyright © 2019-2020 Solus Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	wlog "github.com/DataDrake/waterlog"
)

type status int

const (
	skipped status = iota
	success
	failure
)

type skip struct {
	Chroot    bool     `toml:"chroot,omitempty"`
	LiveMedia bool     `toml:"live_media,omitempty"`
	Paths     []string `toml:"paths"`
}

type check struct {
	Paths []string `toml:"paths"`
}

type replace struct {
	Paths   []string `toml:"paths"`
	Exclude []string `toml:"exclude"`
}

type bin struct {
	Task    string   `toml:"task"`
	Bin     string   `toml:"bin"`
	Args    []string `toml:"args"`
	Replace *replace `toml:"replace"`
}

type remove struct {
	Paths   []string `toml:"paths"`
	Exclude []string `toml:"exclude"`
}

type content struct {
	Description string            `toml:"description"`
	Bins        []*bin            `toml:"bins"`
	Skip        *skip             `toml:"skip,omitempty"`
	Check       *check            `toml:"check,omitempty"`
	Env         map[string]string `toml:"env"`
	RemoveDirs  *remove           `toml:"remove,omitempty"`
}

type output struct {
	Name    string
	SubTask string
	Message string
	Status  status
}

type config struct {
	Name    string
	Path    string
	Output  []*output
	Content *content
}

func load(name string) *config {
	c := &config{
		Name:   name,
		Path:   filepath.Clean(filepath.Join(UsrDir, name+".toml")),
		Output: []*output{{Status: skipped, Name: name}},
	}

	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		c.Path = filepath.Clean(filepath.Join(SysDir, name+".toml"))
		if _, err := os.Stat(c.Path); os.IsNotExist(err) {
			c.Output[0].Status = failure
			c.Output[0].Message = fmt.Sprintf("Unable to find config file located at %s", c.Path)
			return c
		}
	}

	cfg, err := ioutil.ReadFile(c.Path)
	if err != nil {
		c.Output[0].Status = failure
		c.Output[0].Message = fmt.Sprintf("Unable to read config file located at %s", c.Path)
		return c
	}

	cnt := &content{}
	if err := toml.Unmarshal(cfg, cnt); err != nil {
		c.Output[0].Status = failure
		c.Output[0].Message = fmt.Sprintf("Unable to read config file located at %s due to %s", c.Path, err.Error())
		return c
	}
	c.Content = cnt

	bins := c.Content.Bins
	if len(bins) == 0 {
		c.Output[0].Status = failure
		c.Output[0].Message = fmt.Sprintf("%s must contain at least one [[bin]]", c.Name)
		return c
	}

	for _, bin := range bins {
		if len(bin.Task) > 42 {
			c.Output[0].Status = failure
			c.Output[0].Message = fmt.Sprintf("The task: `%s` text cannot exceed the length of 42", bin.Task)
			return c
		}
	}

	return c
}

func (c *config) finish() {
	for _, out := range c.Output {
		t := time.Now()
		now := fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())

		switch out.Status {
		case skipped:
			wlog.Infof("Skipped %s\n", out.Name)
			fmt.Fprintf(os.Stdout, "\033[30;48;5;220m ⁓ \033[7m %s \033[27m %-42s \033[49;38;5;220m %s\033[0m\n", now, out.Name, out.SubTask)
		case failure:
			wlog.Errorf("Failure for %s due to %s\n", out.Name, out.Message)
			fmt.Fprintf(os.Stdout, "\033[30;48;5;208m ✗ \033[7m %s \033[27m %-42s \033[49;38;5;208m %s\033[0m\n", now, out.Name, out.SubTask)
		case success:
			wlog.Goodln(out.Name)
			fmt.Fprintf(os.Stdout, "\033[30;48;5;040m 🗸 \033[7m %s \033[27m %-42s \033[49;38;5;040m %s\033[0m\n", now, out.Name, out.SubTask)
		}
	}
}

func (c *content) skipProcessing() bool {
	check := c.Check
	if check != nil {
		if err := check.checkPaths(); err != nil {
			return true
		}
	}

	if *isForced {
		return false
	}

	skip := c.Skip
	if skip == nil {
		return false
	}

	if skip.Chroot && *isChroot {
		return true
	}

	if skip.LiveMedia && *isLiveMedium {
		return true
	}

	for _, p := range skip.Paths {
		if _, err := os.Stat(filepath.Clean(p)); !os.IsNotExist(err) && !*isForced {
			return true
		}
	}

	return false
}

func (b *bin) execute() error {
	if *isNoRun {
		return nil
	}

	cmd := exec.Command(b.Bin, b.Args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing '%s %v': %s", b.Bin, b.Args, err.Error())
	}

	return nil
}

func (c *check) checkPaths() error {
	for _, path := range c.Paths {
		p, err := filepath.Glob(path)
		if err != nil {
			return fmt.Errorf("unable to glob path")
		}

		if len(p) == 0 {
			return fmt.Errorf("paths are empty")
		}

		for _, pa := range p {
			if _, err := os.Stat(filepath.Clean(pa)); os.IsNotExist(err) {
				return fmt.Errorf("path was not found: %s", pa)
			}
		}
	}

	return nil
}

func (r *remove) removePaths() error {
	paths := filterPaths(r.Paths, r.Exclude)
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			return err
		}
	}

	return nil
}
