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

package cli

import (
    "fmt"
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/waterlog/level"
	"github.com/getsolus/usysconf/config"
	"github.com/getsolus/usysconf/util"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
	wlog "github.com/DataDrake/waterlog"
)

// Run fulfills the "run" subcommand
var Run = cmd.CMD{
    Name:  "run",
    Alias: "r",
    Short: "Run specified configuration file(s) to update the system configuration. It prints the status of each execution: SUCCESS(🗸)/FAILURE(✗)/SKIPPED(⁓)",
    Flags: &RunFlags{},
    Args:  &RunArgs{},
    Run:   RunRun,
}

// RunFlags contains the additional flags for the "run" subcommand
type RunFlags struct {
    Force  bool `short:"f" long:"force"   desc:"Force run the configuration regardless if it should be skipped."`
	DryRun bool `short:"n" long:"dry-run" desc:"Test the configuration files without executing the specified binaries and arguments"`
}

// RunArgs contains the arguments for the "run" subcommand
type RunArgs struct {
    Triggers []string `desc:"Names of the triggers to run"`
}

// RunRun prints the usage for the requested command
func RunRun(r *cmd.RootCMD, c *cmd.CMD) {
    gFlags := r.Flags.(*GlobalFlags)
    args := c.Args.(*RunArgs)
    flags := c.Flags.(*RunFlags)

	// Root user check
	if !flags.DryRun && os.Geteuid() != 0 {
		wlog.Fatalln("You must have root privileges to run triggers")
	}

	// Enable Debug Output
    if gFlags.Debug {
        wlog.SetLevel(level.Debug)
    }

	if !flags.DryRun {
		// Load the system log file
		path := filepath.Clean(filepath.Join(config.LogDir, "usysconf.log"))
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 00600)
		if err != nil {
			wlog.Fatal(err.Error())
		}
		wlog.SetOutput(f)
	}

	wlog.Debugln("Started usysconf")
	defer wlog.Debugln("Exiting usysconf")

	// If the names flag is not present, retrieve the names of the
	// configurations in the system and usr directories.
	n := args.Triggers
	if len(n) == 0 {
		nm := make(map[string]bool)
		n = make([]string, 0)
		ufi, err := ioutil.ReadDir(config.UsrDir)
		if err != nil {
			wlog.Fatalln(err.Error)
		}

		sfi, err := ioutil.ReadDir(config.SysDir)
		if err != nil {
			wlog.Fatalln(err.Error)
		}

		for _, f := range sfi {
			name := strings.Replace(f.Name(), ".toml", "", -1)
			nm[name] = true
			n = append(n, name)
		}

		for _, f := range ufi {
			name := strings.Replace(f.Name(), ".toml", "", -1)
			if _, ok := nm[name]; !ok {
				nm[name] = true
				n = append(n, name)
			}
		}
	}

	for _, name := range n {
		RunConfig(name, flags.Force, gFlags.Chroot, gFlags.Live, flags.DryRun)
	}
}

// RunConfig will process a single configuration.
func RunConfig(name string, isForced, isChroot, isLive, isDryRun bool) {
	cfg := config.Load(name)
	defer cfg.Finish()

	if cfg.Output[0].Status == config.Failure {
		return
	}

	c := cfg.Content

	if c.SkipProcessing(isForced, isChroot, isLive) {
		return
	}

	// Set any environment variables needed to execute the configuratio.
	if err := util.SetEnv(c.Env); err != nil {
		cfg.Output[0].Message = err.Error()
		cfg.Output[0].Status = config.Failure
		return
	}

	rmDirs := c.RemoveDirs
	if rmDirs != nil {
		if err := rmDirs.RemovePaths(); err != nil {
			cfg.Output[0].Message = fmt.Sprintf("error removing path: %s\n", err.Error())
			cfg.Output[0].Status = config.Failure
			return
		}
	}

	bins := c.Bins
	bins, cfg.Output = GetAllBins(bins)

	for i, b := range bins {
		if err := b.Execute(isDryRun); err != nil {
			cfg.Output[i].Message = err.Error()
			cfg.Output[i].Status = config.Failure
			return
		}

		cfg.Output[i].Status = config.Success
	}
}

// GetAllBins Process through the binaries of the configuration and check if
// the "***" replace sequence exists in the arguments and create separate
// binaries to be executed.
func GetAllBins(bins []*config.Bin) ([]*config.Bin, []*config.Output) {
	nbins := make([]*config.Bin, 0)
	outputs := make([]*config.Output, 0)

	for _, b := range bins {
		r := b.Replace

		phExists := false
		phIndex := 0
		for i, arg := range b.Args {
			if r == nil {
				break
			}

			if arg == "***" {
				phExists = true
				phIndex = i
				break
			}
		}

		if !phExists {
			nbins = append(nbins, b)
			out := &config.Output{
				Name:   b.Task,
				Status: config.Skipped,
			}
			outputs = append(outputs, out)
			continue
		}

		wlog.Debugf("replace string exists at arg: %d\n", phIndex)

		paths := util.FilterPaths(r.Paths, r.Exclude)
		for _, p := range paths {
			out := &config.Output{
				Name:    b.Task,
				Status:  config.Skipped,
				SubTask: p,
			}
			b.Args[phIndex] = p
			nbins = append(nbins, b)
			outputs = append(outputs, out)
		}
	}

	return nbins, outputs
}
