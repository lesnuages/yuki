package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/desertbit/grumble"
	"github.com/olekukonko/tablewriter"
)

func printSessions(filter string) error {
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
	sort.Slice(rawData, func(i, j int) bool {
		return rawData[i][6] < rawData[j][6]
	})
	if len(Parser.Sessions) > 0 {
		table.AppendBulk(rawData)
		table.Render()
	}
	return nil
}

func sessions(c *grumble.Context) (err error) {
	var sid uint64
	if len(c.Args) > 0 {
		if sid, err = strconv.ParseUint(c.Args[0], 10, 64); err != nil {
			return err
		}
		if _, ok := Parser.Sessions[sid]; ok {
			Parser.CurrentSession = sid
			c.App.SetPrompt(fmt.Sprintf("yuki[%d]>> ", sid))
		}
	} else {
		return printSessions("")
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
		AllowArgs: true,
		Run:       sessions,
		Completer: completeSessions,
	}
	App.AddCommand(sessions)
}
