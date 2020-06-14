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
	"github.com/DataDrake/cli-ng/cmd"
	wlog "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/level"
	"github.com/getsolus/usysconf/config"
	"github.com/getsolus/usysconf/triggers"
	"os"
	"path/filepath"
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

	// Load Triggers
	tm, err := config.LoadAll()
	if err != nil {
		wlog.Fatalf("Failed to load triggers, reason: %s\n", err.Error())
	}

	// If the names flag is not present, retrieve the names of the
	// configurations in the system and usr directories.
	n := args.Triggers
	if len(n) == 0 {
		for k := range tm {
			n = append(n, k)
		}
	}
	// Establish scope of operations
	s := triggers.Scope{
		Chroot: gFlags.Chroot,
		DryRun: flags.DryRun,
		Forced: flags.Force,
		Live:   gFlags.Live,
	}
	// Run triggers
	triggers.Run(tm, s, n)
}
