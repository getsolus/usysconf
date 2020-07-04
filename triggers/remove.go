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
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/usysconf/state"
	"os"
)

// Remove contains paths to be removed from the system.  Tis supports globbing.
type Remove struct {
	Paths   []string `toml:"paths"`
	Exclude []string `toml:"exclude"`
}

// Remove glob the paths and if it exists it will remove it from the system
func (t *Trigger) Remove(s Scope) bool {
	if s.DryRun {
		log.Debugln("   No Paths will be removed during a dry-run\n")
	}
	if t.RemoveDirs == nil {
		log.Debugln("   No Paths to remove\n")
		return true
	}
	m, err := state.Scan(t.RemoveDirs.Paths)
	if err != nil {
		out := Output{
			Status:  Failure,
			Message: fmt.Sprintf("Failed to remove paths for '%s', reason: %s\n", t.Name, err),
		}
		t.Output = append(t.Output, out)
		return false
	}
	m = m.Exclude(t.RemoveDirs.Exclude)
	for k := range m {
		log.Debugf("    Removing path '%s'\n", k)
		if s.DryRun {
			continue
		}
		if err := os.Remove(k); err != nil {
			out := Output{
				Status:  Failure,
				Message: fmt.Sprintf("Failed to remove paths '%s', reason: %s\n", k, err),
			}
			t.Output = append(t.Output, out)
			return false
		}
	}
	return true
}
