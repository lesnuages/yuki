package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/desertbit/grumble"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/lesnuages/yuki/parser"
)

func writeToPcap(c *grumble.Context) (err error) {
	var (
		file *os.File
		ok   bool
		s    parser.Session
	)
	if s, ok = Parser.Sessions[Parser.CurrentSession]; !ok {
		return fmt.Errorf("You must select a session first")
	}
	if len(c.Args) == 0 {
		return fmt.Errorf("You must provide a filepath.")
	}
	filepath := c.Args[0]
	if file, err = os.Create(filepath); err != nil {
		return err
	}
	writer := pcapgo.NewWriter(file)
	writer.WriteFileHeader((uint32)(len(s.Packets)), layers.LinkTypeEthernet)
	defer file.Close()

	for _, packet := range s.Packets {
		writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
	}
	return nil
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
	App.AddCommand(writePcap)
}
