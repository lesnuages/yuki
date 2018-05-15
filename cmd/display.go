package cmd

import (
	"fmt"
	"os"

	"github.com/desertbit/grumble"
	"github.com/olekukonko/tablewriter"
)

func printSummary(c *grumble.Context) error {
	if s, ok := Parser.Sessions[Parser.CurrentSession]; ok {
		s.Summary()
	}
	return nil
}

func printStats(c *grumble.Context) error {
	if s, ok := Parser.Sessions[Parser.CurrentSession]; ok {
		header := []string{"Number of packets", "Transport Layer", "Source", "Destination", "Source Port", "Destination Port"}
		data := []string{fmt.Sprintf("%d", len(s.Packets)), s.Transport, s.SourceIP, s.DestIP, s.SourcePort, s.DestPort}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		table.Append(data)
		table.Render()
	} else {
		return fmt.Errorf("No session selected.")
	}
	return nil
}

func init() {
	summary := &grumble.Command{
		Name: "summary",
		Help: "Prints all the session packets as a hexdump (starting at layer 1)",
		Run:  printSummary,
	}
	stats := &grumble.Command{
		Name: "stats",
		Help: "Displays statistics about the current session",
		Run:  printStats,
	}
	App.AddCommand(summary)
	App.AddCommand(stats)
}
