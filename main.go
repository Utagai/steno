package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/vhs/pkg/lexer"
	"github.com/charmbracelet/vhs/pkg/parser"
)

func parse(tapeFilePath string) ([]parser.Command, error) {
	tapeBuf, err := os.ReadFile("demo.tape")
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	l := lexer.New(string(tapeBuf))
	p := parser.New(l)
	cmds := p.Parse()

	return cmds, nil
}

func main() {
	cmds, err := parse("demo.tape")
	if err != nil {
		log.Fatalf("could not parse tape: %v", err)
	}

	for _, cmd := range cmds {
		fmt.Printf("Cmd: %q; Args: %q; Options: %q\n", cmd.String(), cmd.Args, cmd.Options)
	}
}
