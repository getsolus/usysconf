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
	// "fmt"
	"github.com/DataDrake/cli-ng/cmd"
	// "github.com/getsolus/usysconf/config"
	"github.com/DataDrake/waterlog/level"
	wlog "github.com/DataDrake/waterlog"
)

// List fulfills the "list" subcommand
var List = cmd.CMD{
    Name:  "list",
    Alias: "ls",
    Short: "List available triggers to run (user-specific)",
    Args:  &ListArgs{},
    Run:   ListRun,
}

// ListArgs contains the arguments for the "list" subcommand
type ListArgs struct {}

// ListRun prints the usage for the requested command
func ListRun(r *cmd.RootCMD, c *cmd.CMD) {
    gFlags := r.Flags.(*GlobalFlags)
    // args := c.Args.(*ListArgs)

	// Enable Debug Output
    if gFlags.Debug {
        wlog.SetLevel(level.Debug)
    }

}
