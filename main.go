package main

import (
	"fmt"
	"io"
	"log/slog"
	"math"

	"github.com/spf13/cobra"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func writeOverlay(stream *ffmpeg.Stream, draw DrawCall, xPos, yPos string, opts Options) *ffmpeg.Stream {
	enableFilter := fmt.Sprintf("between(t,%f,%f)", draw.start, draw.end)
	if draw.toEnd {
		enableFilter = fmt.Sprintf("gte(t,%f)", draw.start)
	}

	args := ffmpeg.KwArgs{"enable": enableFilter, "fontsize": opts.fontSize, "fontcolor": opts.fontColor, "x": xPos, "y": yPos}
	if opts.fontConfig != "" {
		args["fontfile"] = opts.fontConfig
	}

	// NOTE: We are specify 0, 0 as the x and y position for the text overlay but these are overriden by the x and y we specify via the ffmpeg.KwArgs
	// above.
	return stream.Drawtext(draw.text, 0, 0, false, args)
}

type DrawCall struct {
	text  string
	start float64
	end   float64
	toEnd bool
}

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

	draws := make([]DrawCall, 0, len(kpes))
	currentDraw := DrawCall{}
	for i := range kpes {
		kpe := kpes[i]
		logger.Info("processing KPE", slog.String("key", kpe.KeyDisplay), slog.Uint64("when", kpe.WhenMS))
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

	stream := ffmpeg.Input(recordingFile)
	for _, draw := range draws {
		logger.Info("Drawing", slog.Any("draw call", draw))
		stream = writeOverlay(stream, draw, opts.xPosition, opts.yPosition, opts)
	}

	if err := stream.Output(opts.outputFile).OverWriteOutput().Run(); err != nil {
		return fmt.Errorf("could not inject text: %w", err)
	}

	return nil
}

func main() {
	var opts Options

	stenoCmd := &cobra.Command{
		Use:   "steno [tape file] [recording file]",
		Short: "Inject keypress overlay onto a charmebracelet vhs recording",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			if opts.verbose {
				logger = slog.Default()
			}

			if err := run(logger, args[0], args[1], opts); err != nil {
				logger.Error("failed to transcribe tape", err)
			}
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

	stenoCmd.Execute()
}
