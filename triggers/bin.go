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
	"bytes"
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/usysconf/util"
	"os/exec"
)

// Bin contains the details of the binary to be executed.
type Bin struct {
	Task    string   `toml:"task"`
	Bin     string   `toml:"bin"`
	Args    []string `toml:"args"`
	Replace *Replace `toml:"replace"`
}

// Execute the binary from the confuration
func (b *Bin) Execute(s Scope, env map[string]string) Output {
	out := Output{Status: Success}
	// if the norun flag is present do not execute the configuration
	if s.DryRun {
		out.Status = Success
		return out
	}
	// Create command
	cmd := exec.Command(b.Bin, b.Args...)
	// Setup environment
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	// Add buffer for output
	var buff bytes.Buffer
	cmd.Stdout = &buff
	cmd.Stderr = &buff
	// Run the command
	if err := cmd.Run(); err != nil {
		out.Status = Failure
		out.Message = fmt.Sprintf("error executing '%s %v': %s\n%s", b.Bin, b.Args, err.Error(), buff.String())
	}
	return out
}

// FanOut generates one or more bin tasks from a given, as needed
func (b Bin) FanOut() (nbins []Bin, outputs []Output) {

	r := b.Replace

	phExists := false
	phIndex := 0
	for i, arg := range b.Args {
		if r == nil {
			break
		}

		if arg == "***" {
			phExists = true
			phIndex = i
			break
		}
	}

	if !phExists {
		nbins = append(nbins, b)
		out := Output{Name: b.Task}
		outputs = append(outputs, out)
		return
	}

	log.Debugf("    Replace string exists at arg: %d\n", phIndex)

	paths := util.FilterPaths(r.Paths, r.Exclude)
	for _, p := range paths {
		out := Output{
			Name:    b.Task,
			SubTask: p,
		}
		b.Args[phIndex] = p
		nbins = append(nbins, b)
		outputs = append(outputs, out)
	}
	return
}
