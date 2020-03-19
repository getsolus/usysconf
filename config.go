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

// Status indicates if the state of the configuration.
type Status int

const (
	// Skipped - The configuration will not be executed.
	Skipped Status = iota
	// Success - The configuration has been executed, without error.
	Success
	// Failure - The configuration will not be executed, due to error.
	Failure
)

// Skip contains details for when the configuration will not be executed, due
// to existing paths, or possible flags passed.  This supports globbing.
type Skip struct {
	Chroot bool     `toml:"chroot,omitempty"`
	Live   bool     `toml:"live,omitempty"`
	Paths  []string `toml:"paths"`
}

// Check contains paths that must exixt to execute the configuration.  This
// supports globbing.
type Check struct {
	Paths []string `toml:"paths"`
}

// Replace contains details to replace a single argument with a path in the
// executed binary.  This supports globbing.
type Replace struct {
	Paths   []string `toml:"paths"`
	Exclude []string `toml:"exclude"`
}

// Bin contains the details of the binary to be executed.
type Bin struct {
	Task    string   `toml:"task"`
	Bin     string   `toml:"bin"`
	Args    []string `toml:"args"`
	Replace *Replace `toml:"replace"`
}

// Remove contains paths to be removed from the system.  Tis supports globbing.
type Remove struct {
	Paths   []string `toml:"paths"`
	Exclude []string `toml:"exclude"`
}

// Content contains all the details of the configuration file to be executed.
type Content struct {
	Description string            `toml:"description"`
	Bins        []*Bin            `toml:"bins"`
	Skip        *Skip             `toml:"skip,omitempty"`
	Check       *Check            `toml:"check,omitempty"`
	Env         map[string]string `toml:"env"`
	RemoveDirs  *Remove           `toml:"remove,omitempty"`
}

// Output contains the details necessary to output the configuration details
// to the user.
type Output struct {
	Name    string
	SubTask string
	Message string
	Status  Status
}

// Config contains all the information for a configuration to be executed and
// output to the user.
type Config struct {
	Name    string
	Path    string
	Output  []*Output
	Content *Content
}

// Load will check the system and user directories, in that order, for a
// configuration file that has the passed name parameter, without the
// extension and will create a config with the specified valus.
func Load(name string) *Config {
	c := &Config{
		Name:   name,
		Path:   filepath.Clean(filepath.Join(SysDir, name+".toml")),
		Output: []*Output{{Status: Skipped, Name: name}},
	}

	// Check if the configuration exists in the system path, if not check
	// the usr path, otherwise return failure
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		c.Path = filepath.Clean(filepath.Join(UsrDir, name+".toml"))
		if _, err := os.Stat(c.Path); os.IsNotExist(err) {
			c.Output[0].Status = Failure
			c.Output[0].Message = fmt.Sprintf("Unable to find config file located at %s", c.Path)
			return c
		}
	}

	// Read the configuration into the program
	cfg, err := ioutil.ReadFile(c.Path)
	if err != nil {
		c.Output[0].Status = Failure
		c.Output[0].Message = fmt.Sprintf("Unable to read config file located at %s", c.Path)
		return c
	}

	// Save the configuration into the content structure
	cnt := &Content{}
	if err := toml.Unmarshal(cfg, cnt); err != nil {
		c.Output[0].Status = Failure
		c.Output[0].Message = fmt.Sprintf("Unable to read config file located at %s due to %s", c.Path, err.Error())
		return c
	}
	c.Content = cnt

	// Verify that there is at least one binary to execute, otherwise there
	// is no need to continue
	bins := c.Content.Bins
	if len(bins) == 0 {
		c.Output[0].Status = Failure
		c.Output[0].Message = fmt.Sprintf("%s must contain at least one [[bin]]", c.Name)
		return c
	}

	// Verify that the text to be outputed to the user does not exceed 42
	// characters to make the outputed lines uniform in length
	for _, bin := range bins {
		if len(bin.Task) > 42 {
			c.Output[0].Status = Failure
			c.Output[0].Message = fmt.Sprintf("The task: `%s` text cannot exceed the length of 42", bin.Task)
			return c
		}
	}

	return c
}

// Finish is the last function to be exexuted by any configuration to output
// details to the user.
func (c *Config) Finish() {
	for _, out := range c.Output {
		t := time.Now()
		now := t.Format("15:04:05")

		switch out.Status {
		case Skipped:
			wlog.Infof("Skipped %s\n", out.Name)
			fmt.Fprintf(os.Stdout, "\033[30;48;5;220m ⁓ \033[7m %s \033[27m %-42s \033[49;38;5;220m %s\033[0m\n", now, out.Name, out.SubTask)
		case Failure:
			wlog.Errorf("Failure for %s due to %s\n", out.Name, out.Message)
			fmt.Fprintf(os.Stdout, "\033[30;48;5;208m ✗ \033[7m %s \033[27m %-42s \033[49;38;5;208m %s\033[0m\n", now, out.Name, out.SubTask)
		case Success:
			wlog.Goodln(out.Name)
			fmt.Fprintf(os.Stdout, "\033[30;48;5;040m 🗸 \033[7m %s \033[27m %-42s \033[49;38;5;040m %s\033[0m\n", now, out.Name, out.SubTask)
		}
	}
}

// SkipProcessing will process the skip and check elements of the configuration
// and see if it should not be executed.
func (c *Content) SkipProcessing() bool {

	// Check if the paths exist, if not skip
	if c.Check != nil {
		if err := c.Check.CheckPaths(); err != nil {
			wlog.Errorln(err.Error())
			return true
		}
	}

	// Even if the skip element exists, if the force flag is present,
	// continue processing
	if *isForced {
		return false
	}

	if c.Skip == nil {
		return false
	}

	// If the skip element exists and the chroot flag is present, skip
	if c.Skip.Chroot && *isChroot {
		return true
	}

	// If the skip element exists and the live flag is present, skip
	if c.Skip.Live && *isLive {
		return true
	}

	// Process through the skip paths, and if one is present within the
	// system, skip
	for _, p := range c.Skip.Paths {
		if _, err := os.Stat(filepath.Clean(p)); !os.IsNotExist(err) {
			return true
		}
	}

	return false
}

// Execute the binary from the confuration
func (b *Bin) Execute() error {
	// if the norun flag is present do not execute the configuration
	if *isNoRun {
		return nil
	}

	cmd := exec.Command(b.Bin, b.Args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing '%s %v': %s", b.Bin, b.Args, err.Error())
	}

	return nil
}

// CheckPaths will glob the paths and if the path does not exist in the system,
// an error is returned
func (c *Check) CheckPaths() error {
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

// RemovePaths will glob the paths and if it exists it will remove it from the
// system
func (r *Remove) RemovePaths() error {
	paths := FilterPaths(r.Paths, r.Exclude)
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			return err
		}
	}

	return nil
}
