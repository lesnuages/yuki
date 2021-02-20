package cmd

import (
	"github.com/desertbit/grumble"
)

func loadFile(c *grumble.Context) error {
	filePath := c.Args.String("filepath")
	return Parser.Parse(filePath)
}

func init() {
	load := &grumble.Command{
		Name: "load",
		Help: "Loads and parses the provided PCAP/PCAPNG file",
		LongHelp: `Usage: load FILEPATH
Loads the PCAP/PCAPNG file located at FILEPATH and parses it.
		`,
		Run:       loadFile,
		Completer: dirCompleter,
		Args: func(a *grumble.Args) {
			a.String("filepath", "path to a PCAP file")
		},
	}
	App.AddCommand(load)
}
