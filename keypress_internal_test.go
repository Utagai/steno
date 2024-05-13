package main

import (
	"testing"
)

func TestParseKeyPressEvents(t *testing.T) {
	type testCase struct {
		name           string
		tapeFile       string
		noSymbol       bool
		expectedEvents []KeyPressEvent
	}

	tcs := []testCase{
		{
			name:           "empty tape file",
			tapeFile:       "empty.tape",
			noSymbol:       false,
			expectedEvents: []KeyPressEvent{},
		},
		{
			name:     "no special symbols",
			tapeFile: "every_symbol.tape",
			noSymbol: true,
			expectedEvents: []KeyPressEvent{
				{KeyDisplay: "<BACKSPACE>", WhenMS: 0},
				{KeyDisplay: "<DELETE>", WhenMS: 0},
				{KeyDisplay: "<CTRL>+", WhenMS: 0},
				{KeyDisplay: "<CTRL>+", WhenMS: 0},
				{KeyDisplay: "<ALT>+", WhenMS: 0},
				{KeyDisplay: "<DOWN>", WhenMS: 0},
				{KeyDisplay: "<UP>", WhenMS: 0},
				{KeyDisplay: "<RIGHT>", WhenMS: 0},
				{KeyDisplay: "<LEFT>", WhenMS: 0},
				{KeyDisplay: "<PAGEDOWN>", WhenMS: 0},
				{KeyDisplay: "<PAGEUP>", WhenMS: 0},
				{KeyDisplay: "<SPACE>", WhenMS: 0},
				{KeyDisplay: "<TAB>", WhenMS: 50},
				{KeyDisplay: "<ESCAPE>", WhenMS: 50},
				{KeyDisplay: "<ENTER>", WhenMS: 50},
				{KeyDisplay: "<PASTE>", WhenMS: 50},
				{KeyDisplay: "a", WhenMS: 50},
				{KeyDisplay: "b", WhenMS: 50},
				{KeyDisplay: "c", WhenMS: 100},
			},
		},
		{
			name:     "special symbols allowed",
			tapeFile: "every_symbol.tape",
			noSymbol: false,
			expectedEvents: []KeyPressEvent{
				{KeyDisplay: "⌫", WhenMS: 0},
				{KeyDisplay: "␡", WhenMS: 0},
				{KeyDisplay: "C-", WhenMS: 0},
				{KeyDisplay: "C-", WhenMS: 0},
				{KeyDisplay: "⎇-", WhenMS: 0},
				{KeyDisplay: "↓", WhenMS: 0},
				{KeyDisplay: "↑", WhenMS: 0},
				{KeyDisplay: "→", WhenMS: 0},
				{KeyDisplay: "←", WhenMS: 0},
				{KeyDisplay: "⤓", WhenMS: 0},
				{KeyDisplay: "⤒", WhenMS: 0},
				{KeyDisplay: "⎵", WhenMS: 0},
				{KeyDisplay: "⇥", WhenMS: 50},
				{KeyDisplay: "⎋", WhenMS: 50},
				{KeyDisplay: "⏎", WhenMS: 50},
				{KeyDisplay: "C-V", WhenMS: 50},
				{KeyDisplay: "a", WhenMS: 50},
				{KeyDisplay: "b", WhenMS: 50},
				{KeyDisplay: "c", WhenMS: 100},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tapePath := "testtapes/" + tc.tapeFile
			events, err := parseKeyPressEvents(tapePath, tc.noSymbol)
			if err != nil {
				t.Fatalf("Failed to parse key press events: %v", err)
			}
			if len(events) != len(tc.expectedEvents) {
				t.Logf("Got events: %+v", events)
				t.Logf("Expected events: %+v", tc.expectedEvents)
				t.Fatalf("Mismatched number of events. Expected %d, got %d", len(tc.expectedEvents), len(events))
			}
			for i, event := range events {
				if event != tc.expectedEvents[i] {
					t.Fatalf("Mismatched event at index %d. Expected %+v, got %+v", i, tc.expectedEvents[i], event)
				}
			}
		})
	}
}
