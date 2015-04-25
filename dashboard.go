package monitrondashboard

// Basic dashboard functionality

import (
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
)

// rect is a simple struct giving a bounding rectangle for a widget
type rect struct {
	x int
	y int
	w int
	h int
}

type buildState int

const (
	BuildStateFailed buildState = iota
	BuildStateAcknowledged
	BuildStatePassed
	BuildStateUnknown
)

func (bs buildState) BgColour() termbox.Attribute {
	switch bs {
	case BuildStateFailed:
		return termbox.ColorRed
	case BuildStateAcknowledged:
		return termbox.ColorYellow
	case BuildStatePassed:
		return termbox.ColorGreen
	case BuildStateUnknown:
		return termbox.ColorMagenta
	}
	return termbox.ColorCyan
}

func (bs buildState) FgColour() termbox.Attribute {
	switch bs {
	case BuildStateFailed:
		return termbox.ColorWhite
	case BuildStateAcknowledged:
		return termbox.ColorWhite
	case BuildStatePassed:
		return termbox.ColorWhite
	case BuildStateUnknown:
		return termbox.ColorWhite
	}
	return termbox.ColorWhite
}

type build struct {
	name         string
	buildState   buildState
	building     bool
	acknowledger string
}

const textPadding int = 1

// drawBuildState draws a status box for an individual build within bounds
func drawBuildState(build build, bounds rect) {
	bgColour := build.buildState.BgColour()
	fgColour := build.buildState.FgColour()

	availableWidth := bounds.w - 2*textPadding
	buildNameWithLengthRestriction, _ := elipsize(build.name, availableWidth)
	runeName := []rune(buildNameWithLengthRestriction)

	for x := 0; x < bounds.w; x++ {
		for y := 0; y < bounds.h; y++ {
			char := ' '

			// draw build name
			if y == 0 && x >= textPadding {
				charIndex := x - textPadding
				if charIndex < len(runeName) {
					char = runeName[charIndex]
				}
			}

			// TODO show building status

			termbox.SetCell(x+bounds.x, y+bounds.y, char, fgColour, bgColour)
		}
	}
}

// ellipsize returns a string restricted to maxLength using an ellipsis
func elipsize(s string, maxLength int) (string, error) {
	if maxLength < 3 {
		return "", errors.New("Max length too short to ellipsize.")
	}
	runes := []rune(s)
	fmt.Printf("Length: %s => %d\n", s, len(runes))
	if len(runes) > maxLength {
		return fmt.Sprintf("%s...", string(runes[0:maxLength-3])), nil
	}
	return s, nil
}

// redraw redraws the screen.
func redraw() {
	//x, y := termbox.Size()

	drawBuildState(build{"test really really really really really really really really really really really really long", BuildStateFailed, false, ""}, rect{1, 1, 30, 3})
	drawBuildState(build{"test", BuildStatePassed, true, ""}, rect{1, 5, 30, 3})
	drawBuildState(build{"test", BuildStateAcknowledged, false, ""}, rect{1, 9, 30, 3})

	termbox.Flush()
}

func Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	redraw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			default:
				if ev.Ch == 'q' {
					break mainloop
				}
			}
		case termbox.EventError:
			fmt.Printf("Error: %s\n", ev.Err)
			break mainloop
		}
		redraw()
	}
}
