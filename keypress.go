package main

import (
	"time"

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
		true:  "<CTRL>",
		false: "C-",
	},
	token.ALT: {
		true:  "<ALT>",
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
		false: "↵",
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

func parseKeyPressEvents(cmds []parser.Command, noSymbol bool) ([]KeyPressEvent, error) {
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
