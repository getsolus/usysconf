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
	"fmt"
	"log"
	"os"

	wlog "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	flag "github.com/spf13/pflag"
)

var (
	// Version of the application
	Version string
	// UsrDir is the path defined during build i.e. /usr/share/defaults/usysconf.d
	UsrDir string
	// SysDir is the path defined during build i.e. /etc/usysconf.d
	SysDir string
	// LogDir is the path defined during build i.e. /var/log/usysconf
	LogDir string

	isDebug   *bool
	isForced  *bool
	isChroot  *bool
	isLive    *bool
	isNoRun   *bool
	isHelp    *bool
	isVersion *bool
	names     *[]string
)

func init() {
	wlog.SetLevel(level.Info)
	wlog.SetFormat(format.Un)
	wlog.SetFlags(log.Ltime | log.Ldate | log.LUTC)

	isDebug = flag.BoolP("debug", "d", false, "Run in debug mode")
	isForced = flag.BoolP("force", "f", false, "Force run the binaries with the config file(s) regardless if skipped is set")
	isChroot = flag.BoolP("chroot", "c", false, "Specify that command is being run from a chrooted environment")
	isLive = flag.BoolP("live", "l", false, "Specify that command is being run from a live medium")
	isNoRun = flag.Bool("norun", false, "Test the loading of the config file(s) without executing the bin(s)")
	isVersion = flag.BoolP("version", "v", false, "Print the version number of usysconf")
	isHelp = flag.BoolP("help", "h", false, "Print usage information")

	names = flag.StringSliceP("names", "n", []string{}, "Specify the config names to run")

	flag.Parse()

	if *isDebug {
		wlog.SetLevel(level.Debug)
	}
}

func main() {
	args := os.Args
	wlog.Debugf("args: %s\n", args)

	if *isVersion {
		version()
		os.Exit(0)
	}

	if *isHelp {
		usage()
		os.Exit(0)
	}

	if len(args) < 2 {
		wlog.Fatalln("invalid number of arguments")
	}

	switch args[1] {
	case "run":
		run()
	case "version":
		version()
	case "help":
		help(args)
	default:
		wlog.Fatalf("%s is an unknown command\n", args[1])
	}

	os.Exit(0)
}

func version() {
	fmt.Fprintf(os.Stdout, "usysconf v%s\n", Version)
}
