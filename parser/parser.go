package parser

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Parser struct {
	CurrentSession uint64
	Sessions       map[uint64]Session
}

type Session struct {
	Packets    []gopacket.Packet
	Transport  string
	SourceIP   string
	DestIP     string
	SourcePort string
	DestPort   string
	TimeStamp  time.Time
}

func (s *Session) addPacket(p gopacket.Packet) {
	s.Packets = append(s.Packets, p)
}

func (s *Session) Summary() {
	for _, p := range s.Packets {
		fmt.Println(p.Dump())
	}
}

func (p *Parser) createSession(hash uint64) Session {
	p.Sessions[hash] = Session{}
	return p.Sessions[hash]
}

func (p *Parser) Parse(path string) error {
	var (
		ok      bool
		current Session
	)
	if handle, err := pcap.OpenOffline(path); err != nil {
		return err
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			if networkLayer := packet.NetworkLayer(); networkLayer != nil {
				hash := networkLayer.NetworkFlow().FastHash()
				// New session
				if current, ok = p.Sessions[hash]; !ok {
					current = p.createSession(hash)
					current.SourceIP = packet.NetworkLayer().NetworkFlow().Src().String()
					current.DestIP = packet.NetworkLayer().NetworkFlow().Dst().String()
					// Only handle layer 4 datagrams if
					// there is a layer 4
					if packet.TransportLayer() != nil {
						current.SourcePort = packet.TransportLayer().TransportFlow().Src().String()
						current.DestPort = packet.TransportLayer().TransportFlow().Dst().String()
						current.Transport = packet.TransportLayer().LayerType().String()
					}
					current.TimeStamp = packet.Metadata().Timestamp
				}
				// Session already exists
				current.Packets = append(current.Packets, packet)
				p.Sessions[hash] = current
			}
		}
	}
	return nil
}

func NewParser() *Parser {
	return &Parser{
		Sessions: make(map[uint64]Session),
	}
}
