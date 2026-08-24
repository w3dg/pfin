package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/w3dg/pfin/command"
	"github.com/w3dg/pfin/config"
	"github.com/w3dg/pfin/tracker"
)

var (
	configFlag = flag.String("config", "~/.config/pfin/config.toml", "Config file to read from. By default it reads ~/.config/pfin/config.toml . If it doesn't exist, one will be generated.")
	testFlag   = flag.Bool("testmode", false, "If test mode is enabled, we read/write from/to testdata/ during development")
)

func main() {
	flag.Usage = usage

	// Parse stops at first non flag looking arg, so something like "add" contain all the flags (for the add flagset behind it in flag.Args)
	// This way we can have new flags per subcommand
	flag.Parse()

	var cfg config.Config
	var conferr error

	if *testFlag {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not figure working dir for test mode")
		}
		testConfPath := path.Join(cwd, "testdata/testconfig.toml")
		cfg, conferr = loadConfig(&testConfPath)
	} else {
		cfg, conferr = loadConfig(configFlag)
	}

	if conferr != nil {
		log.Fatalf("Not able to parse config: %v", conferr)
	}

	t := tracker.NewTracker(cfg)

	args := flag.Args()
	if len(args) == 0 {
		usage()
	}

	cmdName := args[0]
	cmd, ok := command.Commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "pfin: unknown command %q\n\n", cmdName)
		usage()
	}

	cmd.Dispatch(&t, args[1:])
}

func loadConfig(configFilePath *string) (config.Config, error) {
	if *configFilePath != "~/.config/pfin/config.toml" {
		return config.LoadConfig(*configFilePath)
	}
	return config.ParseConfigOrWriteDefault()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: pfin [global options] <command> [options specific to command]\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	for name, cmd := range command.Commands {
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", name, cmd.Desc)
	}
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}
