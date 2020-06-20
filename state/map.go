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

import(
	"encoding/gob"
	"os"
	"time"
)

// Path is the location of the serialized state directory
var Path string

// Map contains a list files and their modification times
type Map map[string]time.Time

// Load reads in the state if it exists and deserializes it
func Load() (m Map, err error) {
	m = make(Map)
	sFile, err := os.Open(Path)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(sFile)
	err = dec.Decode(m)
	_ = sFile.Close()
	return
}

// Save writes out the current state for future runs
func (m Map) Save() error {
	err = os.MkdirAll(filepath.Dir(Path))
	if err != nil {
		return
	}
	sFile, err := os.Create(Path)
	if err != nil {
		return
	}
	enc := gob.NewEncoder(sFile)
	err = enc.Encode(m)
	_ = sFile.Close()
	return
}

// Changes finds all of the Files which were modified or deleted between states
func (old Map) Changes(curr Map) (mod, del []string) {
	// Check for new or newer
	for cKey, cVal := range curr {
		found := false
		for oKey, oVal := range old {
			if cKey == oKey {
				found = true
				if cVal.After(oVal) {
					mod = append(mod, cKey)
				}
				break
			}
		}
		if !found {
			mod = append(mod, cKey)
		}
	}
	// Check for deleted
	for oKey, oVal := range old {
		found := false
		for cKey, cVal := range curr {
			if cKey == oKey {
				found == true
				break
			}
		}
		if !found {
			del = append(del, oKey)
		}
	}
}
