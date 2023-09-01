// Copyright © Solus Project
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
	"github.com/alecthomas/kong"
)

// Version will be injected by ld flags.
var Version string = "unknown"

// GlobalFlags contains the flags for all commands.
type GlobalFlags struct {
	Debug   bool             `short:"d" long:"debug"  help:"Run in debug mode."`
	Chroot  bool             `short:"c" long:"chroot" help:"Specify that command is being run from a chrooted environment."`
	Live    bool             `short:"l" long:"live"   help:"Specify that command is being run from a live medium."`
	Version kong.VersionFlag `short:"v" long:"version"   help:"Print version and exit."`
}

type arguments struct {
	GlobalFlags

	Run   run   `cmd:"" aliases:"r" help:"Run specified trigger(s) to update the system configuration."`
	List  list  `cmd:"" aliases:"ls" help:"List available triggers to run (user-specific)."`
	Graph graph `cmd:"" aliases:"g" help:"Print the dependencies for all available triggers."`
}

func Parse() (*kong.Context, GlobalFlags) {
	var args arguments
	ctx := kong.Parse(&args, kong.Vars{"version": Version})

	return ctx, args.GlobalFlags
}
