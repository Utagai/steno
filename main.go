package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/vhs/pkg/lexer"
	"github.com/charmbracelet/vhs/pkg/parser"

	ffmpeg "github.com/u2takey/ffmpeg-go"
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

func injectText(text string) error {
	return ffmpeg.Input("demo.mp4").
		Drawtext("hehe", 0, 0, false, ffmpeg.KwArgs{"enable": "between(t,0,3)", "fontsize": 30, "fontcolor": "white", "fontfile": "'./notosans.ttf\\:style=bold'"}).
		Drawtext(text, 0, 0, false, ffmpeg.KwArgs{"enable": "gte(t,3)", "fontsize": 30, "fontcolor": "white", "fontfile": "'./notosans.ttf\\:style=bold'"}).
		Output("output.mp4").OverWriteOutput().ErrorToStdOut().Run()
}

func main() {
	cmds, err := parse("demo.tape")
	if err != nil {
		log.Fatalf("could not parse tape: %v", err)
	}

	for _, cmd := range cmds {
		fmt.Printf("Cmd: %q; Args: %q; Options: %q\n", cmd.String(), cmd.Args, cmd.Options)
	}

	if err := injectText("Hello, world ←!"); err != nil {
		return
	}
}
