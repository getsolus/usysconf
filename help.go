package main

import (
	"fmt"
	"os"

	wlog "github.com/DataDrake/waterlog"
	flag "github.com/spf13/pflag"
)

func help(args []string) {
	helpCommand := flag.NewFlagSet("help", flag.ContinueOnError)

	if err := helpCommand.Parse(args[2:]); err != nil {
		wlog.Fatalln(err.Error())
	}

	if len(args) < 3 {
		usage()
		return
	}

	switch args[2] {
	case "run":
		runUsage()
	default:
		usage()
	}
}

func runUsage() {
	const runHelpText = `Run specified configuration file(s) to update the system configuration.
It prints the status of each execution: SUCCESS(🗸)/FAILURE(✗)/SKIPPED(⁓)	

Usage:
  usysconf run [flags]
  usysconf run --names a,b,c [flags]
	
Flags:
  -c, --chroot                 specify that the configuration is being run in a 
                               chrooted environment
  -d, --debug                  debug output
  -f, --force                  force run the configuration regardless if it 
                               should be skipped
  -l, --live                   specify that the configuration is being run from
                               a live medium
  -n, --names string,string    specify the name of the files to be run in a 
                               comma separated list
      --norun                  test the configuration files without executing
                               the specified binaries and arguments
`

	fmt.Fprintln(os.Stdout, runHelpText)
}

func usage() {
	const helpText = `usysconf is a tool for managing universal system configurations using TOML
based configuration files.	
 	
Usage:	
  usysconf [command]

Flags:
  -h, --help                   help for usysconf
  -v, --version                print the version number of usysconf

Available Commands:
  help       Help about any command
  run        Run the configuration files
  version    Print the version number of usysconf
	
Use "usysconf help [command] for more information about a command.	
`

	fmt.Fprintln(os.Stdout, helpText)
}
