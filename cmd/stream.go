package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/lesnuages/yuki/parser"
)

func stream(c *grumble.Context) error {
	var (
		err     error
		session parser.Session
	)
	if session, err = Parser.GetSession(); err != nil {
		return err
	}

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
	return nil
}

func init() {
	stream := &grumble.Command{
		Name: "stream",
		Help: "Prints the session stream (see flags for encoding)",
		Flags: func(f *grumble.Flags) {
			f.Bool("x", "hex", true, "Print as hexadecimal dump")
			f.Bool("s", "string", false, "Print as ASCII string")
		},
		Run: stream,
	}
	App.AddCommand(stream)
}
