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

package deps

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

// Graph represents the dependencies shared between triggers.
type Graph map[string][]string

// Insert sets the dependencies for a given trigger.
func (g Graph) Insert(name string, deps []string) {
	g[name] = append(g[name], deps...)
}

// Validate checks the graph for any potential issues.
func (g Graph) Validate(triggers []string) {
	g.CheckCircular()
	g.CheckMissing(triggers)
}

// CheckMissing checks for any missign triggers and prints warnings.
func (g Graph) CheckMissing(triggers []string) {
	for name, deps := range g {
		for _, dep := range deps {
			found := false

			for _, trigger := range triggers {
				if dep == trigger {
					found = true
					break
				}
			}

			if !found {
				slog.Warn("Dependency does not exist", "parent", name, "child", dep)
			}
		}
	}
}

// CheckCircular checks for circular dependencies.
func (g Graph) CheckCircular() {
	var visited []string
	for name := range g {
		if found := g.circular(name, visited); len(found) != 0 {
			last := found[len(found)-1]
			for _, next := range found {
				if next == last {
					break
				}

				found = found[1:]
			}
			// TODO: Return an error instead of panicking.
			slog.Error("Circular dependency", "chain", strings.Join(found, " -> "))
			panic("Circular dependency")
		}
	}
}

func (g Graph) circular(name string, visited []string) (found []string) {
	visited = append(visited, name)
	for _, dep := range g[name] {
		for _, v := range visited {
			if dep == v {
				found = append(found, name, dep)
				return
			}
		}

		if len(found) == 0 {
			found = g.circular(dep, visited)
			if len(found) != 0 {
				return append([]string{name}, found...)
			}
		}
	}

	return
}

// prune all references to things not in the list.
func (g Graph) prune(names []string) {
	for k := range g {
		found := false

		for _, name := range names {
			if k == name {
				found = true
				break
			}
		}

		if !found {
			delete(g, k)
		}
	}

	for k, deps := range g {
		var next []string

		for _, dep := range deps {
			found := false

			for _, name := range names {
				if dep == name {
					found = true
					break
				}
			}

			if found {
				next = append(next, dep)
			}
		}

		g[k] = next
	}
}

// traverse performs a breadth-first traversal of a graph.
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

	return order, remaining
}

// Resolve finds the ideal ordering for a list of triggers.
func (g Graph) Resolve(todo []string) (order []string) {
	g.prune(todo)

	var partial []string
	for len(todo) > 0 {
		partial, todo = g.traverse(todo)
		order = append(order, partial...)
	}

	return
}

// Print renders this graph to a "dot" format.
func (g Graph) Print() {
	names := make([]string, 0, len(g))

	for name := range g {
		names = append(names, name)
	}

	sort.Strings(names)
	fmt.Println("digraph {")

	for _, name := range names {
		deps := g[name]
		sort.Strings(deps)

		for _, dep := range deps {
			fmt.Printf("\t\"%s\" -> \"%s\";\n", name, dep)
		}
	}

	fmt.Println("}")
}
