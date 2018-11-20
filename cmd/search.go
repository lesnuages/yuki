package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/lesnuages/yuki/parser"
)

func searchPattern(pattern []byte, format string, s parser.Session) []string {
	var (
		startIdx int
		endIdx   int
		header   string
		before   string
		after    string
		found    string
		result   []string
	)
	for i, p := range s.Packets {
		if p.TransportLayer() != nil {
			if idx := bytes.Index(p.TransportLayer().LayerPayload(), pattern); idx != -1 {
				startIdx = idx
				endIdx = idx + len(pattern) + 30
				if idx-10 > 0 {
					startIdx = idx - 10
				}
				if endIdx > len(p.TransportLayer().LayerPayload()) {
					endIdx -= 30
				}

				header = color.BlueString("[Packet %d]: ", i)
				switch format {
				case "hex":
					before = color.WhiteString("[...] %x", p.TransportLayer().LayerPayload()[startIdx:idx])
					found = color.RedString("%x", p.TransportLayer().LayerPayload()[idx:idx+len(pattern)])
					after = color.WhiteString("%x [...]", p.TransportLayer().LayerPayload()[idx+len(pattern):endIdx])
				case "ascii":
					before = color.WhiteString("[...] %s", p.TransportLayer().LayerPayload()[startIdx:idx])
					found = color.RedString("%s", p.TransportLayer().LayerPayload()[idx:idx+len(pattern)])
					after = color.WhiteString("%s [...]", p.TransportLayer().LayerPayload()[idx+len(pattern):endIdx])
				default:
					before = color.WhiteString("[...] %s", p.TransportLayer().LayerPayload()[startIdx:idx])
					found = color.RedString("%s", p.TransportLayer().LayerPayload()[idx:idx+len(pattern)])
					after = color.WhiteString("%s [...]", p.TransportLayer().LayerPayload()[idx+len(pattern):endIdx])
				}
				result = append(result, fmt.Sprintf("%s%s%s%s", header, before, found, after))
			}
		}
	}
	return result
}

func findPattern(pattern []byte, format string) {
	if Parser.CurrentSession != 0 {
		if _, ok := Parser.Sessions[Parser.CurrentSession]; ok {
			printResults(pattern, format, Parser.CurrentSession)
		}
	} else {
		for h := range Parser.Sessions {
			printResults(pattern, format, h)
		}
	}
}

func printResults(pattern []byte, format string, sid uint64) {
	headline := color.New(color.FgGreen, color.Bold)
	session := Parser.Sessions[sid]
	res := searchPattern(pattern, format, session)
	if len(res) > 0 {
		headline.Printf("[*] Session %d\n", sid)
		for _, r := range res {
			fmt.Println(r)
		}
	}
}

func search(c *grumble.Context) (err error) {
	var (
		pattern []byte
		format  string
	)
	if len(c.Args) < 1 {
		return fmt.Errorf("you must provide a pattern to search")
	}
	arg := c.Args[0]
	if c.Flags.Bool("hex") {
		if pattern, err = hex.DecodeString(arg); err != nil {
			return err
		}
		format = "hex"
	} else {
		pattern = []byte(arg)
		format = "ascii"
	}
	findPattern(pattern, format)
	return nil
}

func init() {
	search := &grumble.Command{
		Name: "search",
		Help: "Look for patterns in session packets",
		LongHelp: `Usage: search [-s | -h] PATTERN
Looks for PATTERN in all sessions packets.
Default behaviour is to look for string patterns (-s implied).
		`,
		Run:       search,
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.Bool("s", "string", true, "ASCII pattern")
			f.Bool("h", "hex", false, "Hex pattern")
		},
	}
	App.AddCommand(search)
}
