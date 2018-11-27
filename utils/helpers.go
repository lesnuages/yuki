package utils

import (
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/lesnuages/yuki/parser"
)

// GetSession returns a parser.Session object given a session id
func GetSession(c *grumble.Context, p *parser.Parser) (parser.Session, error) {
	var (
		ok      bool
		session parser.Session
	)
	if session, ok = p.Sessions[p.CurrentSession]; ok {
		return session, nil
	}
	return session, fmt.Errorf("wrong session selected")
}
