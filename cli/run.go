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
	"errors"
	"fmt"
	"os"

	"github.com/getsolus/usysconf/config"
	"github.com/getsolus/usysconf/triggers"
	"github.com/getsolus/usysconf/util"
)

type run struct {
	Force  bool `short:"f" long:"force"   help:"Force run the configuration regardless if it should be skipped."`
	DryRun bool `short:"n" long:"dry-run" help:"Test the configuration files without executing the specified binaries and arguments."`

	Triggers []string `arg:"" help:"Names of the triggers to run." optional:""`
}

func (r run) Run(flags GlobalFlags) error {
	if os.Geteuid() != 0 {
		return errors.New("you must have root privileges to run triggers")
	}

	if util.IsChroot() {
		flags.Chroot = true
	}

	if util.IsLive() {
		flags.Live = true
	}
	// Load Triggers.
	tm, err := config.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load triggers: %w", err)
	}
	// If the names flag is not present, retrieve the names of the
	// configurations in the system and usr directories.
	n := r.Triggers
	if len(n) == 0 {
		for k := range tm {
			n = append(n, k)
		}
	}
	// Establish scope of operations.
	s := triggers.Scope{
		Chroot: flags.Chroot,
		Debug:  flags.Debug,
		DryRun: r.DryRun,
		Forced: r.Force,
		Live:   flags.Live,
	}
	// Run triggers.
	tm.Run(s, n)

	return nil
}
