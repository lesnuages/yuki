package parser

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Parser describes the parser object.
type Parser struct {
	CurrentSession uint64
	Sessions       map[uint64]Session
	Domains        map[string]string
}

// Session describes the session object.
type Session struct {
	Packets    []gopacket.Packet
	Transport  string
	SourceIP   string
	DestIP     string
	SourcePort string
	DestPort   string
	DomainSrc  string
	DomainDst  string
	TimeStamp  time.Time
}

func (s *Session) addPacket(p gopacket.Packet) {
	s.Packets = append(s.Packets, p)
}

// Summary prints a summary of the current Session.
func (s *Session) Summary() {
	for _, p := range s.Packets {
		fmt.Println(p.Dump())
	}
}

func (p *Parser) createSession(hash uint64) Session {
	p.Sessions[hash] = Session{}
	return p.Sessions[hash]
}

// Parse takes a file path to a PCAP/PCAPNG file
// and extract the content into the Parser.Sessions
// attribute.
func (p *Parser) Parse(path string) error {
	var (
		ok      bool
		current Session
		handle  *pcap.Handle
		err     error
	)
	if handle, err = pcap.OpenOffline(path); err != nil {
		return err
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if networkLayer := packet.NetworkLayer(); networkLayer != nil {
			hash := networkLayer.NetworkFlow().FastHash()
			p.getDNS(packet)
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
			if domain, found := p.Domains[current.SourceIP]; found {
				current.DomainSrc = domain
			}
			if domain, found := p.Domains[current.DestIP]; found {
				current.DomainDst = domain
			}
			p.Sessions[hash] = current
		}

	}
	return nil
}

func (p *Parser) getDNS(packet gopacket.Packet) {
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)
		if dns.QR {
			for _, aa := range dns.Answers {
				if aa.Type.String() == "A" {
					p.Domains[aa.IP.String()] = string(aa.Name[:])
				}
			}
		}
	}
}

// GetSession returns the Session object
// corresponding to Parser.CurrentSession.
// Returns an error if no Session object is found.
func (p *Parser) GetSession() (Session, error) {
	var (
		ok      bool
		session Session
	)
	if session, ok = p.Sessions[p.CurrentSession]; ok {
		return session, nil
	}
	return session, fmt.Errorf("wrong session selected")
}

// NewParser returns a new Parser object.
func NewParser() *Parser {
	return &Parser{
		Sessions: make(map[uint64]Session),
		Domains:  make(map[string]string),
	}
}
