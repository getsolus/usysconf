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
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/usysconf/triggers"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Load reads in all of the trigger files in a directory
func Load(path string) (tm triggers.Map, err error) {
	tm = make(triggers.Map)
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		log.Debugf("Skipped directory '%s':\n", path)
		if os.IsNotExist(err) {
			log.Debugf("    Not found.\n")
			err = nil
			return
		}
		log.Debugf("    Failed to read triggers, reason: %s\n", err)
		err = nil
		return
	}
	log.Debugf("Scanning directory '%s':\n", path)
	found := false
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
		log.Debugf("    Found '%s'\n", t.Name)
		found = true
		if err = t.Load(t.Path); err == nil {
			// Check the config for problems
			err = t.Validate()
		}
		if err != nil {
			err = fmt.Errorf("failed to read '%s' from '%s' reason: %s", name, path, err.Error())
			return
		}
		// save trigger
		tm[t.Name] = t
	}
	if !found {
		log.Debugln("    No triggers found.")
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
	tm.Merge(tm2)
	// Read from Home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	// replace the root directory with the user home directory executing usysconf
	if os.Getuid() == 0 {
		username := os.Getenv("SUDO_USER")
		if username == "" || username == "root" {
			// if user is not found or it is actually being run by root without sudo return
			log.Warnln("Home Triggers not loaded")
			goto CHECK
		}
		// Lookup sudo user's home directory
		u, err := user.Lookup(username)
		if err != nil {
			log.Warnf("Failed to lookup user '%s', reason: %s\n", username, err)
		} else {
			home = u.HomeDir
		}
	}
	// Load configs from the user's Home directory
	tm2, err = Load(filepath.Join(home, ".config", "usysconf.d"))
	if err != nil {
		return
	}
	tm.Merge(tm2)
CHECK:
	// check for lack of triggers
	if len(tm) == 0 {
		log.Fatalln("No triggers available")
	}
	log.Goodf("Found '%d' triggers\n", len(tm))
	return
}
