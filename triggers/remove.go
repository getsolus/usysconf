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

package triggers

import (
	wlog "github.com/DataDrake/waterlog"
	"github.com/getsolus/usysconf/util"
	"os"
)

// Remove contains paths to be removed from the system.  Tis supports globbing.
type Remove struct {
	Paths   []string `toml:"paths"`
	Exclude []string `toml:"exclude"`
}

// Execute will glob the paths and if it exists it will remove it from the system
func (r *Remove) Execute(s Scope) error {
	if s.DryRun {
		wlog.Debugln("No Paths will be removed during a dry-run\n")
	}
	paths := util.FilterPaths(r.Paths, r.Exclude)
	for _, p := range paths {
		wlog.Debugf("Removing path '%s'\n", p)
		if s.DryRun {
			continue
		}
		if err := os.Remove(p); err != nil {
			return err
		}
	}

	return nil
}
