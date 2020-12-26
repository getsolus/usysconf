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
)

// GraphDeps fulfills the "graph-deps" subcommand
var GraphDeps = cmd.CMD{
	Name:  "graph-deps",
	Alias: "dot",
	Short: "Print the dependencies for all available triggers",
	Args:  &GraphDepsArgs{},
	Run:   GraphDepsRun,
}

// GraphDepsArgs contains the arguments for the "graph-deps" subcommand
type GraphDepsArgs struct{}

// GraphDepsRun prints the usage for the requested command
func GraphDepsRun(r *cmd.RootCMD, c *cmd.CMD) {
	gFlags := r.Flags.(*GlobalFlags)
	// args := c.Args.(*GraphDepsArgs)

	// Enable Debug Output
	if gFlags.Debug {
		log.SetLevel(level.Debug)
	}
	// Load Triggers
	tm, err := config.LoadAll()
	if err != nil {
		log.Fatalf("Failed to load triggers, reason: %s\n", err)
	}
	// Print triggers
	tm.Graph(gFlags.Chroot, gFlags.Live)
}
