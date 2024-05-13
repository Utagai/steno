package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/charmbracelet/vhs/pkg/lexer"
	"github.com/charmbracelet/vhs/pkg/parser"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func parseCmds(tapeFilePath string) ([]parser.Command, error) {
	tapeBuf, err := os.ReadFile("demo.tape")
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	l := lexer.New(string(tapeBuf))
	p := parser.New(l)
	cmds := p.Parse()

	return cmds, nil
}

func writeOverlay(stream *ffmpeg.Stream, draw DrawCall) *ffmpeg.Stream {
	enableFilter := fmt.Sprintf("between(t,%f,%f)", draw.start, draw.end)
	if draw.toEnd {
		enableFilter = fmt.Sprintf("gte(t,%f)", draw.start)
	}
	return stream.Drawtext(draw.text, 0, 0, false, ffmpeg.KwArgs{"enable": enableFilter, "fontsize": 30, "fontcolor": "white", "fontfile": "'./notosans.ttf\\:style=bold'"})
}

type DrawCall struct {
	text  string
	start float64
	end   float64
	toEnd bool
}

func main() {
	cmds, err := parseCmds("demo.tape")
	if err != nil {
		log.Fatalf("could not parse tape: %v", err)
	}

	for _, cmd := range cmds {
		fmt.Printf("Cmd: %q; Args: %q; Options: %q\n", cmd.String(), cmd.Args, cmd.Options)
	}

	kpes, err := parseKeyPressEvents(cmds)
	if err != nil {
		log.Fatalf("could not parse key press events: %v", err)
	}

	draws := make([]DrawCall, 0, len(kpes))
	currentDraw := DrawCall{}
	for i := range kpes {
		kpe := kpes[i]
		fmt.Printf("KPE: %q @ %dms\n", kpe.KeyDisplay, kpe.WhenMS)
		currentDraw.text += kpe.KeyDisplay
		if i < len(kpes)-1 {
			currentDraw.end = float64(kpes[i+1].WhenMS) / 1000
		} else {
			currentDraw.end = math.MaxFloat64
			currentDraw.toEnd = true
		}
		currentDraw.start = float64(kpe.WhenMS) / 1000
		draws = append(draws, currentDraw)
	}

	stream := ffmpeg.Input("demo.mp4")
	for _, draw := range draws {
		fmt.Printf("Drawing: %q from %f to %f\n", draw.text, draw.start, draw.end)
		stream = writeOverlay(stream, draw)
	}

	if err := stream.Output("output.mp4").OverWriteOutput().ErrorToStdOut().Run(); err != nil {
		log.Fatalf("could not inject text: %v", err)
	}
}
