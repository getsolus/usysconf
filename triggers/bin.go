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

package triggers

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/getsolus/usysconf/util"
)

// Bin contains the details of the binary to be executed.
type Bin struct {
	Task    string   `toml:"task"`
	Bin     string   `toml:"bin"`
	Args    []string `toml:"args"`
	Replace *Replace `toml:"replace"`
}

// ExecuteBins generates and runs all of the necesarry Bin commands.
func (t *Trigger) ExecuteBins(s Scope) {
	var bins []Bin

	var outputs []Output
	// Generate
	for _, b := range t.Bins {
		bs, outs := b.FanOut()
		bins = append(bins, bs...)
		outputs = append(outputs, outs...)
	}
	// Execute
	for i, b := range bins {
		out := b.Execute(s, t.Env)
		outputs[i].Status = out.Status
		outputs[i].Message = out.Message
	}

	t.Output = append(t.Output, outputs...)
}

// Execute the binary from the confuration.
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

// FanOut generates one or more bin tasks from a given, as needed by replacing the "***" sequence
// in the arguments and creating separate binaries to be executed.
func (b Bin) FanOut() (nbins []Bin, outputs []Output) {
	phIndex := -1

	for i, arg := range b.Args {
		if arg == "***" {
			phIndex = i
			break
		}
	}

	if phIndex == -1 {
		nbins = append(nbins, b)
		out := Output{Name: b.Task}
		outputs = append(outputs, out)

		return
	}

	if b.Replace == nil {
		slog.Error("Placeholder found, but [bins.replaces] is missing")
		return
	}

	slog.Debug("Replace string exists", "argument", phIndex)

	paths := util.FilterPaths(b.Replace.Paths, b.Replace.Exclude)
	for _, path := range paths {
		out := Output{
			Name:    b.Task,
			SubTask: path,
		}
		b.Args[phIndex] = path
		nbins = append(nbins, b)
		outputs = append(outputs, out)
	}

	return
}
