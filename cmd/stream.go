package cmd

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

func stream(c *grumble.Context) error {
	var (
		sessionID uint64
		err       error
	)
	if Parser.CurrentSession != 0 {
		sessionID = Parser.CurrentSession
	} else if len(c.Args) != 0 {
		if sessionID, err = strconv.ParseUint(c.Args[0], 10, 64); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("you need to select a session first")
	}
	if session, ok := Parser.Sessions[sessionID]; ok {
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

func init() {
	stream := &grumble.Command{
		Name: "stream",
		Help: "Prints the session stream (see flags for encoding)",
		Flags: func(f *grumble.Flags) {
			f.Bool("h", "hex", true, "Print as hexadecimal dump")
			f.Bool("s", "string", false, "Print as ASCII string")
		},
		AllowArgs: true,
		Run:       stream,
	}
	App.AddCommand(stream)
}
