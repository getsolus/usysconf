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
	"os"
	"path/filepath"
)

// Check contains paths that must exixt to execute the configuration.  This
// supports globbing.
type Check struct {
	Paths []string `toml:"paths"`
}

// ResolvePaths will glob the paths and if the path does not exist in the system, an error is returned
func (c *Check) ResolvePaths() error {
	for _, path := range c.Paths {
		p, err := filepath.Glob(path)
		if err != nil {
			return fmt.Errorf("unable to glob path")
		}

		if len(p) == 0 {
			return fmt.Errorf("paths are empty")
		}

		for _, pa := range p {
			if _, err := os.Stat(filepath.Clean(pa)); os.IsNotExist(err) {
				return fmt.Errorf("path was not found: %s", pa)
			}
		}
	}

	return nil
}
