package cmd

import (
	"github.com/desertbit/grumble"
	"github.com/lesnuages/yuki/parser"
)

var (
	App = grumble.New(&grumble.Config{
		Name:        "yuki",
		Prompt:      "yuki>> ",
		Description: "Command line packet parser",
		Flags: func(f *grumble.Flags) {
			f.String("f", "filepath", "", "Path to PCAP file.")
		},
	})
	Parser = parser.NewParser()
)

func init() {
	App.OnInit(func(a *grumble.App, flags grumble.FlagMap) (err error) {
		err = nil
		filepath := flags.String("filepath")
		if filepath != "" {
			if err = Parser.Parse(filepath); err != nil {
				return err
			}
		}
		return
	})
}
