package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

func stream(c *grumble.Context) error {
	if Parser.CurrentSession != 0 {
		if session, ok := Parser.Sessions[Parser.CurrentSession]; ok {
			for idx, packet := range session.Packets {
				if len(packet.TransportLayer().LayerPayload()) != 0 {
					header := color.New(color.FgCyan, color.Bold)
					header.Printf("[Packet %d]\n", idx)
					if c.Flags.Bool("string") {
						fmt.Println(string(packet.TransportLayer().LayerPayload()[:]))
					} else {
						fmt.Println(hex.EncodeToString(packet.TransportLayer().LayerPayload()))
					}
				}
			}
		}
		return nil
	}
	return fmt.Errorf("you need to select a session first")
}

func init() {
	stream := &grumble.Command{
		Name: "stream",
		Help: "Prints the session stream (see flags for encoding)",
		Flags: func(f *grumble.Flags) {
			f.Bool("h", "hex", true, "Print as hexadecimal dump")
			f.Bool("s", "string", false, "Print as ASCII string")
		},
		Run: stream,
	}
	App.AddCommand(stream)
}
