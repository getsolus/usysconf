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

package cli

import (
	"fmt"

	"github.com/getsolus/usysconf/config"
	"golang.org/x/exp/slog"
)

type list struct{}

func (l list) Run(flags GlobalFlags) error {
	tm, err := config.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load triggers: %w", err)
	}
	slog.Info("Available triggers:\n\n")
	tm.Print(flags.Chroot, flags.Live)
	return nil
}
