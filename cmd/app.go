package cmd

import (
	"strconv"

	"github.com/desertbit/grumble"
	"github.com/lesnuages/yuki/parser"
)

var (
	// App - Global application variable.
	// Exported to be reachable from every modules
	App = grumble.New(&grumble.Config{
		Name:        "yuki",
		Prompt:      "yuki>> ",
		Description: "Command line packet parser",
		Flags: func(f *grumble.Flags) {
			f.String("f", "filepath", "", "Path to PCAP file.")
			f.String("s", "sid", "", "Session ID, required by certain commands")
		},
	})
	// Parser - Global parser to be used by all commands.
	Parser = parser.NewParser()
)

func init() {
	App.OnInit(func(a *grumble.App, flags grumble.FlagMap) (err error) {
		err = nil
		filepath := flags.String("filepath")
		sessionID := flags.String("sid")
		if filepath != "" {
			if err = Parser.Parse(filepath); err != nil {
				return err
			}
		}
		if sessionID != "" {
			if sid, err := strconv.ParseUint(sessionID, 10, 64); err == nil {
				Parser.CurrentSession = sid
			} else {
				return err
			}
		}
		return
	})
}
