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
)

// Check contains paths that must exixt to execute the configuration.
// This supports globbing.
type Check struct {
	Paths []string `toml:"paths"`
}

// CheckMatch will glob the paths and if the path does not exist in the system, an error is returned
func (t *Trigger) CheckMatch() (m state.Map, ok bool) {
	ok = true
	if t.Check == nil {
		log.Debugf("No check paths for trigger '%s'\n", t.Name)
		return
	}
	m, err := state.Scan(t.Check.Paths)
	if err != nil {
		out := Output{
			Status:  Failure,
			Message: fmt.Sprintf("Failed to scan paths for '%s', reason: %s\n", t.Name, err),
		}
		t.Output = append(t.Output, out)
		ok = false
		return
	}
	return
}
