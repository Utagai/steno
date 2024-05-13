package main

import (
	"fmt"
	"log/slog"
	"math"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type DrawCall struct {
	text  string
	start float64
	end   float64
	toEnd bool
}

func drawCallsFromKeyPressEvents(logger *slog.Logger, kpes []KeyPressEvent) []DrawCall {
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

	return draws
}

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

func transcribeToVideo(logger *slog.Logger, kpes []KeyPressEvent, recordingFile string, opts Options) error {
	draws := drawCallsFromKeyPressEvents(logger, kpes)

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
