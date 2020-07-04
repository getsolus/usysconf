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
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/usysconf/state"
)

// Trigger contains all the information for a configuration to be executed and
// output to the user.
type Trigger struct {
	Name   string
	Path   string
	Output []Output

	Description string            `toml:"description"`
	Bins        []Bin             `toml:"bins"`
	Skip        *Skip             `toml:"skip,omitempty"`
	Check       *Check            `toml:"check,omitempty"`
	Env         map[string]string `toml:"env"`
	RemoveDirs  *Remove           `toml:"remove,omitempty"`
}

// Run will process a single configuration and scope.
func (t *Trigger) Run(s Scope, prev, next state.Map) (ok bool) {
	var check, diff state.Map
	// Get the new check result
	check, ok = t.CheckMatch()
	if !ok {
		goto FINISH
	}
	// Calculate Diff
	diff = state.Diff(prev, check)
	// Merge it into the new State
	next.Merge(diff)
	// Check for Skip
	if t.ShouldSkip(s, check, diff) {
		goto FINISH
	}
	// Do the removals
	if ok = t.Remove(s); !ok {
		goto FINISH
	}
	// Run the bins
	t.ExecuteBins(s)
FINISH:
	t.Finish(s)
	return
}

// Finish is the last function to be executed by any trigger to output details to the user.
func (t *Trigger) Finish(s Scope) {
	// Check for the worst status
	status := Skipped
	for _, out := range t.Output {
		if out.Status > status {
			status = out.Status
		}
	}
	// Indicate the worst status for the whole group
	switch status {
	case Skipped:
		log.Debugln(t.Name)
	case Failure:
		log.Errorln(t.Name)
	case Success:
		log.Goodln(t.Name)
	}
	// Indicate status for sub-tasks
	for _, out := range t.Output {
		switch out.Status {
		case Skipped:
			if len(out.SubTask) > 0 {
				log.Debugf("    Skipped for %s due to %s\n", out.SubTask, out.Message)
			} else if len(out.Message) > 0 {
				log.Debugf("    Skipped due to %s\n", out.Message)
			}
		case Failure:
			if len(out.SubTask) > 0 {
				log.Errorf("    Failure for %s due to %s\n", out.SubTask, out.Message)
			} else if len(out.Message) > 0 {
				log.Errorf("    Failure due to %s\n", out.Message)
			}
		case Success:
			if s.DryRun && len(out.SubTask) > 0 {
				log.Infof("    %s\n", out.SubTask)
			}
		}
	}
}
