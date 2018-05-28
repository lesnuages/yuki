package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/lesnuages/yuki/parser"
)

func dumpToFile(c *grumble.Context, isPcap bool) (err error) {
	var (
		file *os.File
		ok   bool
		s    parser.Session
	)
	if s, ok = Parser.Sessions[Parser.CurrentSession]; !ok {
		return fmt.Errorf("You must select a session first")
	}
	if len(c.Args) == 0 {
		return fmt.Errorf("You must provide a filepath")
	}
	filepath := c.Args[0]
	if file, err = os.Create(filepath); err != nil {
		return err
	}
	defer file.Close()
	if isPcap {
		writer := pcapgo.NewWriter(file)
		writer.WriteFileHeader((uint32)(len(s.Packets)), layers.LinkTypeEthernet)
		for _, p := range s.Packets {
			writer.WritePacket(p.Metadata().CaptureInfo, p.Data())
		}
	} else {
		for _, p := range s.Packets {
			file.Write(p.TransportLayer().LayerPayload())
		}
	}
	headline := color.New(color.FgGreen, color.Bold)
	headline.Println("[*] File written to", filepath)
	return nil
}

func writeToPcap(c *grumble.Context) (err error) {
	return dumpToFile(c, true)
}

func writeToFile(c *grumble.Context) (err error) {
	return dumpToFile(c, false)
}

func dirCompleter(prefix string, args []string) []string {
	var (
		res        []string
		currentDir string
	)
	idx := strings.LastIndexAny(prefix, "/")
	if idx == -1 {
		currentDir = "./"
	} else {
		if len(prefix) > 0 {
			currentDir = prefix[:idx+1]
		}
	}
	if files, err := ioutil.ReadDir(currentDir); err == nil {
		for _, file := range files {
			if strings.HasPrefix(currentDir+file.Name(), prefix) {
				index := strings.LastIndex(currentDir+file.Name(), prefix)
				if index == -1 {
					res = append(res, file.Name())
				} else if index < len(currentDir+file.Name()) {
					filename := currentDir + file.Name()[index:]
					res = append(res, filename)
				}
			}
		}
	} else {
		fmt.Println(err)
	}
	return res
}

func init() {
	writePcap := &grumble.Command{
		Name:      "writepcap",
		Help:      "Export the current session to a pcap file",
		AllowArgs: true,
		Run:       writeToPcap,
		Completer: dirCompleter,
	}
	dump := &grumble.Command{
		Name:      "dump",
		Help:      "Dump transport layer data to a file",
		AllowArgs: true,
		Run:       writeToFile,
		Completer: dirCompleter,
	}
	App.AddCommand(writePcap)
	App.AddCommand(dump)
}
