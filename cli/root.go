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
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"log"
)

// GlobalFlags contains the flags for all commands
type GlobalFlags struct{
	Debug  bool `short:"d" long:"debug"  desc:"Run in debug mode"`
	Chroot bool `short:"c" long:"chroot" desc:"Specify that command is being run from a chrooted environment"`
	Live   bool `short:"l" long:"live"   desc:"Specify that command is being run from a live medium"`
}

// Root is the main command for this application
var Root *cmd.RootCMD

func init() {
	// Build Application
	Root = &cmd.RootCMD{
		Name:  "usysconf",
		Short: "A tool for managing universal system configurations using TOML based configuration files",
		Flags: &GlobalFlags{},
	}
	// Setup the Sub-Commands
	Root.RegisterCMD(&cmd.Help)
	Root.RegisterCMD(&Run)
	Root.RegisterCMD(&List)
	Root.RegisterCMD(&Version)

	//Set up logging
	wlog.SetLevel(level.Info)
    wlog.SetFormat(format.Un)
    wlog.SetFlags(log.Ltime | log.Ldate | log.LUTC)
}
