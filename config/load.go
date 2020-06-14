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

package config

import (
	"fmt"
	wlog "github.com/DataDrake/waterlog"
	"github.com/getsolus/usysconf/triggers"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Load reads in all of the trigger files in a directory
func Load(path string) (tm triggers.Map, err error) {
	tm = make(triggers.Map)
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		wlog.Debugf("Failed to read triggers in %s, reason: %s\n", path, err.Error())
		err = nil
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".toml") {
			continue
		}
		t := triggers.Trigger{
			Name: strings.TrimSuffix(name, ".toml"),
			Path: filepath.Clean(filepath.Join(path, name)),
		}
		// found trigger
		wlog.Debugf("Found trigger '%s' in directory '%s'\n", name, path)
		if err = t.Config.Load(t.Path); err == nil {
			// Check the config for problems
			err = t.Config.Validate()
		}
		if err != nil {
			err = fmt.Errorf("failed to read '%s' from '%s' reason: %s", name, path, err.Error())
			return
		}
		// save trigger
		tm[t.Name] = t
	}
	return
}

// LoadAll will check the system, user, and home directories, in that order, for a
// configuration file that has the passed name parameter, without the extension
// and will create a config with the specified valus.
func LoadAll() (tm triggers.Map, err error) {
	// Read from System directory
	tm, err = Load(SysDir)
	if err != nil {
		return
	}

	// Read from User Directory
	tm2, err := Load(UsrDir)
	if err != nil {
		return
	}
	triggers.Merge(tm, tm2)

	// Read from Home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	tm2, err = Load(filepath.Join(home, ".config", "usysconf.d"))
	if err != nil {
		return
	}
	triggers.Merge(tm, tm2)
	return
}
