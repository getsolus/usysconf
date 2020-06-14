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
	"sort"
)

// Map relates the name of trigger to its definition
type Map map[string]Trigger

// Merge combines two Maps by copying from right to left
func Merge(left, right Map) {
	for k, v := range right {
		left[k] = v
	}
}

// Print renders a Map in a human-readable format
func Print(tm Map) {
	var keys []string
	max := 0
	for k := range tm {
		keys = append(keys, k)
		if len(k) > max {
			max = len(k)
		}
	}
	sort.Strings(keys)
	var t Trigger
	f := fmt.Sprintf("\t%%%ds - %%s\n", max)
	for _, key := range keys {
		t = tm[key]
		fmt.Printf(f, t.Name, t.Config.Description)
	}
	fmt.Println()
}

// Run executes a list of triggers, where available
func Run(tm Map, s Scope, names []string) {
	// Iterate over triggers
	for _, name := range names {
		// Get Trigger if available
		t, ok := tm[name]
		if !ok {
			wlog.Warnf("Could not find trigger %s\n", name)
			continue
		}
		// Run Trigger
		t.Run(s)
	}
}
