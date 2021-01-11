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
	"github.com/getsolus/usysconf/deps"
	"github.com/getsolus/usysconf/state"
	"sort"
)

// Map relates the name of trigger to its definition
type Map map[string]Trigger

// Merge combines two Maps by copying from right to left
func (tm Map) Merge(tm2 Map) {
	for k, v := range tm2 {
		tm[k] = v
	}
}

// Print renders a Map in a human-readable format
func (tm Map) Print(chroot, live bool) {
	var keys []string
	max := 0
	for k := range tm {
		keys = append(keys, k)
		if len(k) > max {
			max = len(k)
		}
	}
	max += 4
	sort.Strings(keys)
	f := fmt.Sprintf("%%%ds - %%s\n", max)
	for _, key := range keys {
		t := tm[key]
		if t.Skip != nil {
			if (t.Skip.Chroot && chroot) || (t.Skip.Live && live) {
				continue
			}
		}
		log.Printf(f, t.Name, t.Description)
	}
	log.Println()
}

// Graph generates a dependency graph
func (tm Map) Graph(chroot, live bool) (g deps.Graph) {
	g = make(deps.Graph)
	for _, t := range tm {
		if t.Skip != nil {
			if (t.Skip.Chroot && chroot) || (t.Skip.Live && live) {
				continue
			}
		}
		if t.Deps != nil {
			g.Insert(t.Name, t.Deps.After)
		}
	}
	return
}

// Run executes a list of triggers, where available
func (tm Map) Run(s Scope, names []string) {
	prev := state.Load()
	next := make(state.Map)
	// Resolve deps
	g := tm.Graph(s.Chroot, s.Live)
	order := g.Resolve(names)
	// Iterate over triggers
	for _, name := range order {
		// Get Trigger if available
		t, ok := tm[name]
		if !ok {
			log.Warnf("Could not find trigger %s\n", name)
			continue
		}
		// Run Trigger
		t.Run(s, prev, next)
	}
	if !s.DryRun {
		// Save new State for next run
		if err := next.Save(); err != nil {
			log.Errorf("Failed to save next state file, reason: %s\n", err)
		}
	}
}
