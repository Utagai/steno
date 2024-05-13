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
var CommandTypes = []CommandType{ //nolint: deadcode
	token.SET,
		For this one we need to track:
		- Typing speed
		- Playback speed
	token.SOURCE,
		Yea... gonna need some recursion + flattening?
}
*/

func parseKeyPressEvents(cmds []parser.Command) ([]KeyPressEvent, error) {
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
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "⌫", WhenMS: t})
		case token.DELETE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "␡", WhenMS: t})
		case token.CTRL:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "C-", WhenMS: t})
		case token.ALT:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "⎇-", WhenMS: t})
		case token.DOWN:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "↓", WhenMS: t})
		case token.PAGEDOWN:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "⤓", WhenMS: t})
		case token.UP:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "↑", WhenMS: t})
		case token.PAGEUP:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "⤒", WhenMS: t})
		case token.LEFT:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "←", WhenMS: t})
		case token.RIGHT:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "→", WhenMS: t})
		case token.SPACE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "⎵", WhenMS: t})
		case token.ENTER:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "↵", WhenMS: t})
		case token.ESCAPE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "⎋", WhenMS: t})
		case token.TAB:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "⇥", WhenMS: t})
		case token.PASTE:
			keyPressEvents = append(keyPressEvents, KeyPressEvent{KeyDisplay: "C-V", WhenMS: t})
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
