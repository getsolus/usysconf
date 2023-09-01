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

package config

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/getsolus/usysconf/triggers"
)

// Load reads in all of the trigger files in a directory.
func Load(path string) (triggers.Map, error) {
	logger := slog.With("path", path)

	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug("Directory not found")
			return nil, nil
		}

		return nil, fmt.Errorf("failed to read triggers: %w", err)
	}

	tm := make(triggers.Map, len(entries))

	logger.Debug("Scanning directory")

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
		logger.Debug("Trigger found", "name", t.Name)

		err = t.Load(t.Path)
		if err != nil {
			return nil, err
		}

		err = t.Validate()
		if err != nil {
			return nil, fmt.Errorf("failed to read %s from %s: %w", name, path, err)
		}

		tm[t.Name] = t
	}

	if len(tm) == 0 {
		logger.Debug("No triggers found")
	}

	return tm, nil
}

// LoadAll will check the system, user, and home directories, in that order, for a
// configuration file that has the passed name parameter, without the extension
// and will create a config with the specified valus.
func LoadAll() (triggers.Map, error) {
	paths := []string{SysDir, UsrDir}
	if p, err := os.UserHomeDir(); err != nil {
		paths = append(paths, p)
	}

	if os.Getuid() == 0 {
		uname := os.Getenv("SUDO_USER")
		if uname != "" && uname != "root" {
			u, err := user.Lookup(uname)
			if err != nil {
				slog.Warn("Failed to lookup underlying user", "name", uname, "reason", err)
			} else {
				paths = append(paths, filepath.Join(u.HomeDir, ".config", "usysconf.d"))
			}
		}
	}

	tm := make(triggers.Map)

	for _, path := range paths {
		trig, err := Load(path)
		if err != nil {
			return nil, err
		}

		tm.Merge(trig)
	}

	slog.Info("Total triggers", "count", len(tm))

	return tm, nil
}
