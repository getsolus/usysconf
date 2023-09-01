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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Load reads a Trigger configuration from a file and parses it
func (t *Trigger) Load(path string) error {
	// Check if this is a valid file path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	// Read the configuration into the program
	cfg, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("unable to read config file located at %s", path)
	}
	// Save the configuration into the content structure
	if err := toml.Unmarshal(cfg, t); err != nil {
		return fmt.Errorf("unable to read config file located at %s due to %s", path, err.Error())
	}
	return nil
}

// Validate checks for errors in a Trigger configuration
func (t *Trigger) Validate() error {
	// Verify that there is at least one binary to execute, otherwise there
	// is no need to continue
	if len(t.Bins) == 0 {
		return fmt.Errorf("triggers must contain at least one [[bin]]")
	}
	return nil
}
