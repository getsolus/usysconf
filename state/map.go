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

package state

import (
	"fmt"
	log "github.com/DataDrake/waterlog"
	cbor "github.com/fxamacker/cbor/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Path is the location of the serialized system state directory
var Path string

// Map contains a list files and their modification times
type Map map[string]time.Time

// Load reads in the state if it exists and deserializes it
func Load() Map {
	m := make(Map)
	sFile, err := os.Open(filepath.Clean(Path))
	if err != nil {
		return m
	}
	dec := cbor.NewDecoder(sFile)
	_ = dec.Decode(m)
	_ = sFile.Close()
	return m
}

// Save writes out the current state for future runs
func (m Map) Save() error {
	if err := os.MkdirAll(filepath.Dir(Path), 0750); err != nil {
		return err
	}
	sFile, err := os.Create(filepath.Clean(Path))
	if err != nil {
		return err
	}
	enc := cbor.NewEncoder(sFile)
	err = enc.Encode(m)
	_ = sFile.Close()
	return err
}

// Merge combines two Maps into one
func (m Map) Merge(other Map) {
	for k, v := range other {
		m[k] = v
	}
}

// Diff finds all of the Files which were modified or deleted between states
func Diff(old, curr Map) Map {
	diff := make(Map)
	// Check for new or newer
	for cKey, cVal := range curr {
		found := false
		for oKey, oVal := range old {
			if cKey == oKey {
				found = true
				if cVal.After(oVal) {
					diff[cKey] = cVal
				}
				break
			}
		}
		if !found {
			diff[cKey] = cVal
		}
	}
	return diff
}

// Search finds all of the matching files in a Map
func (m Map) Search(paths []string) Map {
	match := make(Map)
	for _, path := range paths {
		search := path
		search = strings.ReplaceAll(search, "*", ".*")
		search = "^" + strings.ReplaceAll(search, string(filepath.Separator), "\\"+string(filepath.Separator))
		regex, err := regexp.Compile(search)
		if err != nil {
			log.Warnf("Could not convert path to regex: %s\n", path)
			continue
		}
		for k, v := range m {
			if regex.MatchString(k) {
				match[k] = v
			}
		}
	}
	return match
}

// Exclude removes keys from the Map if they match certain patterns
func (m Map) Exclude(patterns []string) Map {
	match := make(Map)
	var regexes []*regexp.Regexp
	for _, pattern := range patterns {
		exclude := pattern
		exclude = strings.ReplaceAll(exclude, "*", ".*")
		regex, err := regexp.Compile(exclude)
		if err != nil {
			log.Warnf("Could not convert pattern to regex: %s\n", pattern)
			continue
		}
		regexes = append(regexes, regex)
	}
	for k := range m {
		for _, regex := range regexes {
			if regex.MatchString(k) {
				delete(m, k)
				break
			}
		}
	}
	return match
}

// IsEmpty checkes if the Map has nothing in it
func (m Map) IsEmpty() bool {
	for k := range m {
		if len(k) > 0 {
			return false
		}
		return false
	}
	return true
}

// Strings gets a list of files from the keys
func (m Map) Strings() []string {
	var strs []string
	for k := range m {
		strs = append(strs, k)
	}
	return strs
}

// Scan goes over a set of paths and imports them and their contents to the map
func Scan(paths []string) (m Map, err error) {
	m = make(Map)
	var p []string
	for _, path := range paths {
		p, err = filepath.Glob(path)
		if err != nil {
			err = fmt.Errorf("unable to glob path: %s", path)
			return
		}

		if len(p) == 0 {
			continue
		}
		var info os.FileInfo
		for _, pa := range p {
			info, err = os.Stat(filepath.Clean(pa))
			if err != nil && !os.IsNotExist(err) {
				err = fmt.Errorf("failed to check path: %s", pa)
				return
			}
			m[pa] = info.ModTime()
		}
	}
	return
}
