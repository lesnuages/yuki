package cmd

import (
	"fmt"

	"github.com/desertbit/grumble"
)

func loadFile(c *grumble.Context) error {
	if len(c.Args) == 0 {
		return fmt.Errorf("you must provide a filepath")
	}
	filePath := c.Args[0]
	Parser.Parse(filePath)
	return nil
}

func init() {
	load := &grumble.Command{
		Name: "load",
		Help: "Loads and parses the provided PCAP/PCAPNG file",
		LongHelp: `Usage: load FILEPATH
Loads the PCAP/PCAPNG file located at FILEPATH and parses it.
		`,
		AllowArgs: true,
		Run:       loadFile,
		Completer: dirCompleter,
	}
	App.AddCommand(load)
}
