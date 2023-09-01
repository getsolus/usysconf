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

package main

import (
	log2 "log"

	log "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/getsolus/usysconf/cli"
)

func main() {
	ctx, flags := cli.Parse()

	if flags.Debug {
		log.SetLevel(level.Debug)
	}
	log.SetFormat(format.Min)
	log.SetFlags(log2.Ltime | log2.Ldate | log2.LUTC)

	err := ctx.Run(flags)
	if err != nil {
		log.Fatal(err)
	}
}
