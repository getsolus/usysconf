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

	"github.com/getsolus/usysconf/state"
)

// Skip contains details for when the configuration will not be executed, due to existing paths, or possible flags passed.
// This supports globbing.
type Skip struct {
	Chroot bool     `toml:"chroot,omitempty"`
	Live   bool     `toml:"live,omitempty"`
	Paths  []string `toml:"paths"`
}

// ShouldSkip will process the skip and check elements of the configuration and see if it should not be executed.
func (t *Trigger) ShouldSkip(s Scope, check, diff state.Map) bool {
	out := Output{
		Status: Skipped,
	}
	// Check if the paths exist, if not skip
	if check.IsEmpty() || diff.IsEmpty() {
		t.Output = append(t.Output, out)
		return true
	}
	// Even if the skip element exists, if the force flag is present, continue processing
	if s.Forced {
		return false
	}
	if t.Skip == nil {
		return false
	}
	// If the skip element exists and the chroot flag is present, skip
	if t.Skip.Chroot && s.Chroot {
		t.Output = append(t.Output, out)
		return true
	}
	// If the skip element exists and the live flag is present, skip
	if t.Skip.Live && s.Live {
		t.Output = append(t.Output, out)
		return true
	}
	// Process through the skip paths, and if one is present within the system, skip
	matches := check.Search(t.Skip.Paths)
	for k := range matches {
		out.Message = fmt.Sprintf("path '%s' found", k)
		t.Output = append(t.Output, out)
		return true
	}
	return false
}
