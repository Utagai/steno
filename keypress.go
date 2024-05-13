package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/vhs/pkg/lexer"
	"github.com/charmbracelet/vhs/pkg/parser"
	"github.com/charmbracelet/vhs/pkg/token"
)

type KeyPressEvent struct {
	KeyDisplay string
	// WhenMS is the time in milliseconds when the key was pressed starting from the beginning of the video (t=0).
	WhenMS uint64
}

/**
TODO:
var CommandTypes = []CommandType{ //nolint: deadcode
	token.SET,
		For this one we need to track:
		- Typing speed
		- Playback speed
	token.SOURCE,
		Yea... gonna need some recursion + flattening?
}
*/

// keypressSymbols maps token types to their corresponding keypress description
// string or symbol. The description string and symbol are embedded into an
// inner map, which can be indexed into based on whether special symbols are
// requested or not.
var keypressSymbols = map[token.Type]map[bool]string{
	token.BACKSPACE: {
		true:  "<BACKSPACE>",
		false: "⌫",
	},
	token.DELETE: {
		true:  "<DELETE>",
		false: "␡",
	},
	token.CTRL: {
		true:  "<CTRL>+",
		false: "C-",
	},
	token.ALT: {
		true:  "<ALT>+",
		false: "⎇-",
	},
	token.DOWN: {
		true:  "<DOWN>",
		false: "↓",
	},
	token.PAGEDOWN: {
		true:  "<PAGEDOWN>",
		false: "⤓",
	},
	token.UP: {
		true:  "<UP>",
		false: "↑",
	},
	token.PAGEUP: {
		true:  "<PAGEUP>",
		false: "⤒",
	},
	token.LEFT: {
		true:  "<LEFT>",
		false: "←",
	},
	token.RIGHT: {
		true:  "<RIGHT>",
		false: "→",
	},
	token.SPACE: {
		true:  "<SPACE>",
		false: "⎵",
	},
	token.ENTER: {
		true:  "<ENTER>",
		false: "⏎",
	},
	token.ESCAPE: {
		true:  "<ESCAPE>",
		false: "⎋",
	},
	token.TAB: {
		true:  "<TAB>",
		false: "⇥",
	},
	token.PASTE: {
		true:  "<PASTE>",
		false: "C-V",
	},
}

// NOTE: This function is implemented based on mostly just my observations.
// There are some aspects of it that I think are likely incorrect. For example,
// I only apply the typing speed increment on certain key presses, namely Type
// keys and Space. In my testing, this seems to perform correctly but the only
// way to be absolutely certain is to read through vhs' source code more
// carefully. I could do this, but to be honest, that's a lot of work for a
// project that exists namely to just help another project, and an understanding
// gained this way is brittle and prone to break when vhs internals change.
// This does however suggest that perhaps the true best way to implement the
// purpose of this program is to build it into vhs, but introducing that
// complexity (and dependencies) to vhs is not really desirable either.
// :shrug:
func parseKeyPressEvents(tapeFile string, noSymbol bool) ([]KeyPressEvent, error) {
	cmds, err := parseCmds(tapeFile)
	if err != nil {
		return nil, fmt.Errorf("could not parse tape: %w", err)
	}

	keyPressEvents := make([]KeyPressEvent, 0, len(cmds))
	// t tracks the current time (in milliseconds) in the current
	// demo.
	var t uint64 = 0
	var typeSpeedMs uint64 = 50 // Default.
	// var playbackSpeed float64 = 1.0 // Default.
	for i := range cmds {
		// TODO: I think type events can include a time too, but we have no idea what that means yet and therefore how to handle
		// it.
		cmd := cmds[i]
		switch cmd.Type {
		case token.BACKSPACE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.BACKSPACE][noSymbol], WhenMS: t})
		case token.DELETE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.DELETE][noSymbol], WhenMS: t})
		case token.CTRL:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.CTRL][noSymbol], WhenMS: t})
		case token.ALT:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.ALT][noSymbol], WhenMS: t})
		case token.DOWN:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.DOWN][noSymbol], WhenMS: t})
		case token.PAGEDOWN:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.PAGEDOWN][noSymbol], WhenMS: t})
		case token.UP:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.UP][noSymbol], WhenMS: t})
		case token.PAGEUP:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.PAGEUP][noSymbol], WhenMS: t})
		case token.LEFT:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.LEFT][noSymbol], WhenMS: t})
		case token.RIGHT:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.RIGHT][noSymbol], WhenMS: t})
		case token.SPACE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.SPACE][noSymbol], WhenMS: t})
			t += typeSpeedMs
		case token.ENTER:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.ENTER][noSymbol], WhenMS: t})
		case token.ESCAPE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.ESCAPE][noSymbol], WhenMS: t})
		case token.TAB:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.TAB][noSymbol], WhenMS: t})
		case token.PASTE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: keypressSymbols[token.PASTE][noSymbol], WhenMS: t})
		case token.TYPE:
			for _, r := range cmd.Args {
				keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: string(r), WhenMS: t})
				t += typeSpeedMs
			}
			// Subtract away typeSpeedMs one more time at the end, since there is no delay _after_ the last character is typed.
			// NOTE: I _think_ vhs doesn't introduce an initial delay before typing the first character either. If they did,
			// then I think the simpler thing is to move the increment in the loop to _before_ the append() call, but I have
			// found the current code to produce the most accurate transcription of keypresses.
			t -= typeSpeedMs
		case token.SLEEP:
			tArg, err := time.ParseDuration(cmd.Args)
			if err != nil {
				return nil, err
			}
			t += uint64(tArg.Milliseconds())
		}
	}

	return keyPressEvents, nil
}

func parseCmds(tapeFile string) ([]parser.Command, error) {
	tapeBuf, err := os.ReadFile(tapeFile)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	l := lexer.New(string(tapeBuf))
	p := parser.New(l)
	cmds := p.Parse()

	return cmds, nil
}
