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

package util

import (
	"path/filepath"
)

// get a list of files that match the provided filters
func match(filters []string) (matches []string) {
	for _, filter := range filters {
		partial, err := filepath.Glob(filter)
		if err != nil {
			continue
		}

		matches = append(matches, partial...)
	}

	return
}

// FilterPaths will process through globbed paths and remove any paths from the resulting slice if they are present in the excludes slice.
func FilterPaths(includes []string, excludes []string) (paths []string) {
	excludePaths := match(excludes)
	for _, includePath := range match(includes) {
		for _, excludePath := range excludePaths {
			if includePath == excludePath {
				break
			}
		}

		paths = append(paths, includePath)
	}

	return
}
