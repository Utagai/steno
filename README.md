# `steno`

`steno` is a Go program that is meant to be used with [`vhs`](https://github.com/charmbracelet/vhs) tapes and recording videos. The idea is that `steno` acts as a stenographer for your `vhs` recordings. You give it your `.tape` file and the output `.mp4` video, and it will output a video file with your keypresses overlayed on the video. See the two examples below to get a better sense of what it is that `steno` does:

| Original `vhs` Recording                                                                      | `steno` Post-Processed Recording                                                                |
| --------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| ![demo](https://github.com/Utagai/steno/assets/10730394/66a636c3-9241-490a-8c32-995a2034d0ec) | ![output](https://github.com/Utagai/steno/assets/10730394/cf06ab95-76ac-43af-8f56-e97e6bdf85c7) |

The above comparison demonstrations were made using my project, [`jam`](https://github.com/utagai/jam) and its `vhs` recording!

## Warnings

There's actually plenty of things `steno` is currently doing incorrectly. For example, it isn't handling `PlaybackSpeed` at all, nor is it honoring changes to `TypingSpeed`. It also isn't capturing modifier characters yet. These all should be fixable in due time. There are some other issues too, like handling repeated `Type` events (e.g. `Type@3 4`). And of course, there may be other issues I am not even yet aware of!

## Usage

```
Inject keypress overlay onto a charmbracelet vhs recording

Usage:
  steno [tape file] [recording file] [flags]

Flags:
  -c, --color string         Font color for ffmpeg (default "white")
  -f, --fontconfig string    A fontconfig string for specifying the font and styling for ffmpeg
  -h, --help                 help for steno
      --no-special-symbols   Enable special symbols for keypresses
  -o, --output string        Output file to write to (default "output.mp4")
  -s, --size int             Font size for ffmpeg (default 30)
  -v, --verbose              Enable logging
  -x, --x string             X position for text overlay for ffmpeg; defaults to center (default "(w-text_w)/2")
  -y, --y string             Y position for text overlay for ffmpeg; defaults to bottom (with padding) (default "h-text_h-20")
```
