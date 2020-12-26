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
	log "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/level"
	"github.com/getsolus/usysconf/config"
	"github.com/getsolus/usysconf/triggers"
	"github.com/getsolus/usysconf/util"
	"os"
)

// Run fulfills the "run" subcommand
var Run = cmd.CMD{
	Name:  "run",
	Alias: "r",
	Short: "Run specified trigger(s) to update the system configuration.",
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
	// Enable Debug Output
	if gFlags.Debug {
		log.SetLevel(level.Debug)
	}
	log.Debugln("Started usysconf")
	defer log.Debugln("Exiting usysconf")
	// Root user check
	if os.Geteuid() != 0 {
		log.Fatalln("You must have root privileges to run triggers")
	}
	// Set Chroot as needed
	if util.IsChroot() {
		gFlags.Chroot = true
	}
	// Set Live as needed
	if util.IsLive() {
		gFlags.Live = true
	}
	// Load Triggers
	tm, err := config.LoadAll()
	if err != nil {
		log.Fatalf("Failed to load triggers, reason: %s\n", err)
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
		Debug:  gFlags.Debug,
		DryRun: flags.DryRun,
		Forced: flags.Force,
		Live:   gFlags.Live,
	}
	// Run triggers
	tm.Run(s, n)
}
