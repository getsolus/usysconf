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
	// "fmt"
	wlog "github.com/DataDrake/waterlog"
	// "os"
	//"time"
)

// Trigger contains all the information for a configuration to be executed and
// output to the user.
type Trigger struct {
	Name   string
	Path   string
	Output []Output
	Config Config
}

// Run will process a single configuration and scope.
func (c *Trigger) Run(s Scope) {
	c.Output = c.Config.Execute(s)
	c.Finish(s)
}

// Finish is the last function to be executed by any trigger to output details to the user.
func (c *Trigger) Finish(s Scope) {
	//ansiYellow := "\033[30;48;5;220m"
	//ansiGreen := "\033[30;48;5;040m"
	//ansiRed := "\033[30;48;5;208m"
	//ansiInverse := "\033[7m"
	//ansiInverseReset := "\033[27m"
	//ansiReset := "\033[0m"

	// Check for the worst status
	status := Skipped
	for _, out := range c.Output {
		if out.Status > status {
			status = out.Status
		}
	}
	// Indicate the worst status for the whole group
	switch status {
	case Skipped:
		wlog.Debugln(c.Name)
	case Failure:
		wlog.Errorln(c.Name)
	case Success:
		wlog.Goodln(c.Name)
	}
	// Indicate status for sub-tasks
	for _, out := range c.Output {
		switch out.Status {
		case Skipped:
			if len(out.SubTask) > 0 {
				wlog.Debugf("    Skipped for %s due to %s\n", out.SubTask, out.Message)
			} else if len(out.Message) > 0 {
				wlog.Debugf("    Skipped due to %s\n", out.Message)
			}
		case Failure:
			if len(out.SubTask) > 0 {
				wlog.Errorf("    Failure for %s due to %s\n", out.SubTask, out.Message)
			} else if len(out.Message) > 0 {
				wlog.Errorf("    Failure due to %s\n", out.Message)
			}
		case Success:
			if s.DryRun && len(out.SubTask) > 0 {
				wlog.Infof("    %s\n", out.SubTask)
			}
		}
	}
}
