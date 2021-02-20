package cmd

import (
	"fmt"
	"os"

	"github.com/desertbit/grumble"
	"golang.org/x/crypto/ssh/terminal"
)

func term(c *grumble.Context) error {
	interp := terminal.NewTerminal(os.Stdin, ">")
	if interp == nil {
		return fmt.Errorf("could not create terminal")
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)
	return nil
}

func init() {
	termCmd := &grumble.Command{
		Name: "interpreter",
		Help: "Start an interpreter",
		Run:  term,
	}
	App.AddCommand(termCmd)
}
