// Copyright © 2019-Present Solus Project
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
	"github.com/getsolus/usysconf/state"
	"golang.org/x/exp/slog"
)

// Trigger contains all the information for a configuration to be executed and output to the user.
type Trigger struct {
	Name   string
	Path   string
	Output []Output

	Description string            `toml:"description"`
	Check       *Check            `toml:"check,omitempty"`
	Skip        *Skip             `toml:"skip,omitempty"`
	Deps        *Deps             `toml:"deps,omitempty"`
	Env         map[string]string `toml:"env,omitempty"`
	Bins        []Bin             `toml:"bins,omitempty"`
	Removals    []Remove          `toml:"remove,omitempty"`
}

// Run will process a single configuration and scope.
func (t *Trigger) Run(s Scope, prev, next state.Map) (ok bool) {
	var check, diff state.Map
	// Get the new check result
	if check, ok = t.CheckMatch(); !ok {
		goto FINISH
	}
	// Calculate Diff
	diff = prev.Diff(check)
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
		slog.Debug(t.Name)
	case Failure:
		slog.Error(t.Name)
	case Success:
		slog.Info(t.Name)
	}
	// Indicate status for sub-tasks
	for _, out := range t.Output {
		switch out.Status {
		case Skipped:
			if len(out.SubTask) > 0 {
				slog.Debug("Skipped", "subtask", out.SubTask, "reason", out.Message)
			} else if len(out.Message) > 0 {
				slog.Debug("Skipped", "reason", out.Message)
			}
		case Failure:
			if len(out.SubTask) > 0 {
				slog.Error("Failed", "subtask", out.SubTask, "reason", out.Message)
			} else if len(out.Message) > 0 {
				slog.Error("Failed", "reason", out.Message)
			}
		case Success:
			if s.DryRun && len(out.SubTask) > 0 {
				slog.Info(out.SubTask)
			}
		}
	}
}
