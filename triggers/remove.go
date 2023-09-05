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

package triggers

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/getsolus/usysconf/state"
)

// Remove contains paths to be removed from the system. This supports globbing.
type Remove struct {
	Paths   []string `toml:"paths"`
	Exclude []string `toml:"exclude"`
}

// Remove glob the paths and if it exists it will remove it from the system
func (t *Trigger) Remove(s Scope) bool {
	if s.DryRun {
		slog.Debug("No Paths will be removed during a dry-run")
	}
	if len(t.Removals) == 0 {
		slog.Debug("No Paths to remove")
		return true
	}
	for _, remove := range t.Removals {
		if !t.removeOne(s, remove) {
			return false
		}
	}
	return true
}

// removeOne carries out removals for a single Remove entry
func (t *Trigger) removeOne(s Scope, remove Remove) bool {
	matches, err := state.Scan(remove.Paths)
	if err != nil {
		out := Output{
			Status:  Failure,
			Message: fmt.Sprintf("Failed to remove paths for '%s', reason: %s\n", t.Name, err),
		}
		t.Output = append(t.Output, out)
		return false
	}
	matches = matches.Exclude(remove.Exclude)
	for path := range matches {
		slog.Debug("Removing", "path", path)
		if s.DryRun {
			continue
		}
		if err := os.Remove(path); err != nil {
			out := Output{
				Status:  Failure,
				Message: fmt.Sprintf("Failed to remove path '%s', reason: %s\n", path, err),
			}
			t.Output = append(t.Output, out)
			return false
		}
	}
	return true
}
