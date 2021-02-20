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
		fmtStr   string
		result   []string
	)
	for i, p := range s.Packets {
		if p.TransportLayer() != nil {
			if idx := bytes.Index(p.TransportLayer().LayerPayload(), pattern); idx != -1 {
				if startIdx = idx; idx-10 > 0 {
					startIdx = idx - 10
				}
				if endIdx = idx + len(pattern) + 30; endIdx > len(p.TransportLayer().LayerPayload()) {
					endIdx -= 30
				}

				header = color.BlueString("[Packet %d]: ", i)
				payloadBefore := p.TransportLayer().LayerPayload()[startIdx:idx]
				payloadAfter := p.TransportLayer().LayerPayload()[idx+len(pattern) : endIdx]
				switch format {
				case "hex":
					fmtStr = "%x"
				case "ascii":
					fmtStr = "%s"
				}
				before = color.WhiteString(fmt.Sprintf("[...] %s", fmtStr), payloadBefore)
				found = color.RedString(fmt.Sprintf("%s", fmtStr), pattern)
				after = color.WhiteString(fmt.Sprintf("%s [...]", fmtStr), payloadAfter)
				result = append(result, fmt.Sprintf("%s%s%s%s", header, before, found, after))
			}
		}
	}
	return result
}

func findPattern(pattern []byte, format string) {
	if _, err := Parser.GetSession(); err == nil {
		printResults(pattern, format, Parser.CurrentSession)
	} else {
		for sid := range Parser.Sessions {
			printResults(pattern, format, sid)
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
	arg := c.Args.String("pattern")
	if arg == "" {
		return fmt.Errorf("you must provide a pattern to search")
	}
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
		LongHelp: `Usage: search [-s | -x] PATTERN
Looks for PATTERN in all sessions packets.
Default behaviour is to look for string patterns (-s implied).
		`,
		Run: search,
		Flags: func(f *grumble.Flags) {
			f.Bool("s", "string", true, "ASCII pattern")
			f.Bool("x", "hex", false, "Hex pattern")
		},
		Args: func(a *grumble.Args) {
			a.String("pattern", "pattern to search for")
		},
	}
	App.AddCommand(search)
}
