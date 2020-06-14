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
	wlog "github.com/DataDrake/waterlog"
	"os"
	"time"
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
	c.Finish()
}

// Finish is the last function to be executed by any trigger to output details to the user.
func (c *Trigger) Finish() {
	ansiYellow := "\033[30;48;5;220m"
	ansiGreen := "\033[30;48;5;040m"
	ansiRed := "\033[30;48;5;208m"
	ansiInverse := "\033[7m"
	ansiInverseReset := "\033[27m"
	ansiReset := "\033[0m"

	for _, out := range c.Output {
		t := time.Now()
		now := t.Format("15:04:05")
		name := fmt.Sprintf(" %-42s ", c.Name)

		switch out.Status {
		case Skipped:
			wlog.Warnf("Skipped %s:%s\n", name, out.SubTask)
			fmt.Fprintln(os.Stdout, ansiYellow+" 🗲 "+ansiInverse+" "+now+" "+ansiInverseReset+name+" "+ansiInverse+" "+out.SubTask+ansiReset)
		case Failure:
			wlog.Errorf("Failure for %s:%s due to %s\n", name, out.SubTask, out.Message)
			fmt.Fprintln(os.Stdout, ansiRed+" ✗ "+ansiInverse+" "+now+" "+ansiInverseReset+name+" "+ansiInverse+" "+out.SubTask+ansiReset)
		case Success:
			wlog.Goodf("Succeeded to run %s:%s\n", name, out.SubTask)
			fmt.Fprintln(os.Stdout, ansiGreen+" 🗸 "+ansiInverse+" "+now+" "+ansiInverseReset+name+" "+ansiInverse+" "+out.SubTask+ansiReset)
		}
	}
}
