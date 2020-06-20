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
)

// VersionNumber is the version of the compiled usysconf binary (Makefile)
var VersionNumber string

// Version fulfills the "version" subcommand
var Version = cmd.CMD{
	Name:  "version",
	Alias: "v",
	Short: "Get the version number",
	Args:  &VersionArgs{},
	Run:   VersionRun,
}

// VersionArgs contains the arguments for the "version" subcommand
type VersionArgs struct{}

// VersionRun prints the version of Usysconf and exits
func VersionRun(r *cmd.RootCMD, c *cmd.CMD) {
	// gFlags := r.Flags.(*GlobalFlags)
	// args := c.Args.(*VersionArgs)
	wlog.Infof("usysconf %s\n", VersionNumber)
}
