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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	wlog "github.com/DataDrake/waterlog"
)

func run() {
	if os.Geteuid() != 0 {
		wlog.Fatalln("You must be root to run triggers")
	}

	path := filepath.Clean(filepath.Join(LogDir, "usysconf.log"))
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 00600)
	if err != nil {
		wlog.Fatal(err.Error())
	}
	wlog.SetOutput(f)

	wlog.Debugln("Started usysconf")
	defer wlog.Debugln("Exiting usysconf")

	n := *names
	if len(n) == 0 {
		nm := make(map[string]bool)
		n = make([]string, 0)
		ufi, err := ioutil.ReadDir(UsrDir)
		if err != nil {
			wlog.Fatalln(err.Error)
		}

		sfi, err := ioutil.ReadDir(SysDir)
		if err != nil {
			wlog.Fatalln(err.Error)
		}

		for _, f := range ufi {
			name := strings.Replace(f.Name(), ".toml", "", -1)
			nm[name] = true
			n = append(n, name)
		}

		for _, f := range sfi {
			name := strings.Replace(f.Name(), ".toml", "", -1)
			if _, ok := nm[name]; !ok {
				nm[name] = true
				n = append(n, name)
			}
		}
	}

	for _, name := range n {
		runConfig(name)
	}
}

func runConfig(name string) {
	cfg := load(name)
	defer cfg.finish()

	if cfg.Output[0].Status == failure {
		return
	}

	c := cfg.Content

	if c.skipProcessing() {
		return
	}

	if err := setEnv(c.Env); err != nil {
		cfg.Output[0].Message = err.Error()
		cfg.Output[0].Status = failure
		return
	}

	rmDirs := c.RemoveDirs
	if rmDirs != nil {
		if err := rmDirs.removePaths(); err != nil {
			cfg.Output[0].Message = fmt.Sprintf("error removing path: %s\n", err.Error())
			cfg.Output[0].Status = failure
			return
		}
	}

	bins := c.Bins
	bins, cfg.Output = getAllBins(bins)

	for i, b := range bins {
		if err := b.execute(); err != nil {
			cfg.Output[i].Message = err.Error()
			cfg.Output[i].Status = failure
			return
		}

		cfg.Output[i].Status = success
	}
}

func getAllBins(bins []*bin) ([]*bin, []*output) {
	nbins := make([]*bin, 0)
	outputs := make([]*output, 0)

	for _, b := range bins {
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
			out := &output{
				Name:   b.Task,
				Status: skipped,
			}
			outputs = append(outputs, out)
			continue
		}

		wlog.Debugf("replace string exists at arg: %d\n", phIndex)

		paths := filterPaths(r.Paths, r.Exclude)
		for _, p := range paths {
			out := &output{
				Name:    b.Task,
				Status:  skipped,
				SubTask: p,
			}
			b.Args[phIndex] = p
			nbins = append(nbins, b)
			outputs = append(outputs, out)
		}
	}

	return nbins, outputs
}