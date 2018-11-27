package utils

import (
	"fmt"
	"strconv"

	"github.com/desertbit/grumble"
	"github.com/lesnuages/yuki/parser"
)

// GetSession returns a parser.Session object given a session id
func GetSession(c *grumble.Context, p *parser.Parser) (parser.Session, error) {
	var (
		sid     uint64
		err     error
		ok      bool
		session parser.Session
	)
	if len(c.Args) != 0 {
		if sid, err = strconv.ParseUint(c.Args[0], 10, 64); err != nil {
			return session, err
		}
	} else {
		sid = p.CurrentSession
	}
	if session, ok = p.Sessions[sid]; ok {
		return session, nil
	}
	return session, fmt.Errorf("no session selected")
}
