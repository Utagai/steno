package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	fontConfig string
	fontSize   int
	fontColor  string
	outputFile string
	xPosition  string
	yPosition  string
	noSymbols  bool
	verbose    bool
}

func run(logger *slog.Logger, tapeFile, recordingFile string, opts Options) error {
	kpes, err := parseKeyPressEvents(tapeFile, opts.noSymbols)
	if err != nil {
		return fmt.Errorf("could not parse key press events: %w", err)
	}

	return transcribeToVideo(logger, kpes, recordingFile, opts)
}

func main() {
	var opts Options

	stenoCmd := &cobra.Command{
		Use:   "steno [tape file] [recording file]",
		Short: "Inject keypress overlay onto a charmbracelet vhs recording",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			if opts.verbose {
				logger = slog.Default()
			}

			if err := run(logger, args[0], args[1], opts); err != nil {
				logger.Error("failed to transcribe tape", err)
				return err
			}

			return nil
		},
	}

	stenoCmd.Flags().StringVarP(&opts.fontConfig, "fontconfig", "f", "", "A fontconfig string for specifying the font and styling for ffmpeg")
	stenoCmd.Flags().IntVarP(&opts.fontSize, "size", "s", 30, "Font size for ffmpeg")
	stenoCmd.Flags().StringVarP(&opts.fontColor, "color", "c", "white", "Font color for ffmpeg")
	stenoCmd.Flags().StringVarP(&opts.outputFile, "output", "o", "output.mp4", "Output file to write to")
	stenoCmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Enable logging")
	stenoCmd.Flags().BoolVar(&opts.noSymbols, "no-special-symbols", false, "Enable special symbols for keypresses")
	// NOTE: The default value calculates the horizontal center.
	// We use text_w in the calculation here because we must account for the width of the text we are
	// positioning.
	stenoCmd.Flags().StringVarP(&opts.xPosition, "x", "x", "(w-text_w)/2", "X position for text overlay for ffmpeg; defaults to center")
	// NOTE: The default value calculates an area just above the bottom of the video.
	// We use text_h in the calculation here because we must account for the height of the text we are
	// positioning.
	stenoCmd.Flags().StringVarP(&opts.yPosition, "y", "y", "h-text_h-20", "Y position for text overlay for ffmpeg; defaults to bottom (with padding)")

	if err := stenoCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
