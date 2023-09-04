// Copyright © 2019-Present Solus Project
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
	"os"

	"golang.org/x/exp/slog"
)

// IsLive checks is this process is running in a Live install
func IsLive() bool {
	var err error
	if _, err = os.Stat("/run/initramfs/livedev"); err == nil {
		slog.Debug("Live session detected")
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	slog.Error("Could not check for live session", "reason", err)
	// TODO: Return error instead of panicking.
	panic("Could not check for live session")
	//return false
}
