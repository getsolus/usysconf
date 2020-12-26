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

package deps

import (
	log "github.com/DataDrake/waterlog"
	"sort"
)

// Graph represents the dependencies shared between triggers
type Graph map[string][]string

// Insert sets the dependencies for a given trigger
func (g Graph) Insert(name string, deps []string) {
	g[name] = append(g[name], deps...)
}

// traverse performs a breadth-first traversal of a graph
func (g Graph) traverse(todo []string) (order, remaining []string) {
	for _, name := range todo {
		deps := g[name]
		if len(deps) == 0 {
			order = append(order, name)
		} else {
			remaining = append(remaining, name)
		}
	}
	for _, name := range remaining {
		var next []string
		for _, dep := range g[name] {
			found := false
			for _, prev := range order {
				if dep == prev {
					found = true
					break
				}
			}
			if !found {
				next = append(next, dep)
			}
		}
		g[name] = next
	}
	sort.Strings(order)
	return
}

// Resolve finds the ideal ordering for a list of triggers
func (g Graph) Resolve(todo []string) (order []string) {
	var partial []string
	for len(todo) > 0 {
		partial, todo = g.traverse(todo)
		order = append(order, partial...)
	}
	return
}

// Print renders this graph to a "dot" format
func (g Graph) Print() {
	var names []string
	for name := range g {
		names = append(names, name)
	}
	sort.Strings(names)
	log.Println("digraph {")
	for _, name := range names {
		deps := g[name]
		sort.Strings(deps)
		for _, dep := range deps {
			log.Printf("\t\"%s\" -> \"%s\";\n", name, dep)
		}
	}
	log.Println("}")
}
