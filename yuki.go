package main

import (
	"fmt"

	"github.com/lesnuages/yuki/parser"
)

func main() {
	p := parser.NewParser()
	p.Parse("test.pcap")
	for hash, session := range p.Sessions {
		fmt.Println(hash, len(session.Packets))
	}
}
