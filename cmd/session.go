package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/desertbit/grumble"
	"github.com/olekukonko/tablewriter"
)

func printSessions(filter string) {
	var rawData [][]string
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Session ID",
		"Source IP",
		"Source Port",
		"Destination IP",
		"Destination Port",
		"Source Domain",
		"Destination Domain",
		"Transport Type",
		"Timestamp",
	}
	table.SetHeader(headers)
	for hash, s := range Parser.Sessions {
		content := []string{fmt.Sprint(hash),
			s.SourceIP,
			s.SourcePort,
			s.DestIP,
			s.DestPort,
			s.DomainSrc,
			s.DomainDst,
			s.Transport,
			s.TimeStamp.Format(time.ANSIC),
		}
		rawData = append(rawData, content)
	}
	// Sort by timestamp DESC
	sort.Slice(rawData, func(i, j int) bool {
		return rawData[i][6] < rawData[j][6]
	})
	if len(Parser.Sessions) > 0 {
		table.AppendBulk(rawData)
		table.Render()
	}
}

func sessions(c *grumble.Context) (err error) {
	sid := c.Flags.Uint64("session-id")
	if sid != 0 {
		if _, ok := Parser.Sessions[sid]; ok {
			Parser.CurrentSession = sid
			c.App.SetPrompt(fmt.Sprintf("yuki[%d]>> ", sid))
		}
	} else {
		printSessions("")
	}
	return nil
}

func completeSessions(prefix string, args []string) []string {
	res := make([]string, 0, len(Parser.Sessions))
	for k := range Parser.Sessions {
		strKey := fmt.Sprintf("%d", k)
		if strings.HasPrefix(strKey, prefix) {
			res = append(res, strKey)
		}
	}
	return res
}

func init() {
	sessions := &grumble.Command{
		Name:      "sessions",
		Help:      "List sessions",
		Run:       sessions,
		Completer: completeSessions,
		Flags: func(f *grumble.Flags) {
			f.Uint64("s", "session-id", 0, "session identifier")
		},
	}
	App.AddCommand(sessions)
}
