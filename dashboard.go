package monitrondashboard

// Basic dashboard functionality

import (
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
)

const OrangeColour int = 167

// rect is a simple struct giving a bounding rectangle for a widget
type rect struct {
	x int
	y int
	w int
	h int
}

// rect defines a single point on the screen
type point struct {
	x int
	y int
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
		return termbox.Attribute(OrangeColour)
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

var buildingMessage string = "Building"

// drawBuildState draws a status box for an individual build within bounds
func drawBuildState(build build, bounds rect) {
	bgColour := build.buildState.BgColour()
	fgColour := build.buildState.FgColour()

	availableWidth := bounds.w - 2*textPadding
	buildNameWithLengthRestriction, _ := elipsize(build.name, availableWidth)

	nameWriter := createTextWriter(buildNameWithLengthRestriction, point{1, 0})

	var buildingWriter func(char rune, point point) rune
	var usernameWriter func(char rune, point point) rune
	if build.building {
		buildingWriter = createTextWriter(buildingMessage, point{1, 1})
		usernameWriter = createTextWriter(build.acknowledger, point{10, 1})
	} else {
		buildingWriter = func(char rune, point point) rune { return char }
		usernameWriter = createTextWriter(build.acknowledger, point{1, 1})
	}

	for x := 0; x < bounds.w; x++ {
		for y := 0; y < bounds.h; y++ {
			char := ' '
			currentPoint := point{x, y}
			char = nameWriter(char, currentPoint)
			char = buildingWriter(char, currentPoint)
			char = usernameWriter(char, currentPoint)

			termbox.SetCell(x+bounds.x, y+bounds.y, char, fgColour, bgColour)
		}
	}
}

// createTextPrinter takes a string and a starting point and returns a function,
// the returned function takes the current char to be displayed and a point and
// will return the the character this printer thinks should be displayed at this point.
// Note these functions take the existing char to enable multiple functions to be
// chained together
func createTextWriter(text string, startingPoint point) func(char rune, point point) rune {
	runeText := []rune(text)
	expectedRow := startingPoint.y
	columnOffset := startingPoint.x
	return func(char rune, point point) rune {
		if point.y == expectedRow && point.x >= columnOffset {
			charIndex := point.x - columnOffset
			if charIndex < len(runeText) {
				return runeText[charIndex]
			}
		}
		return char
	}
}

// ellipsize returns a string restricted to maxLength using an ellipsis
func elipsize(s string, maxLength int) (string, error) {
	if maxLength < 3 {
		return "", errors.New("Max length too short to ellipsize.")
	}
	runes := []rune(s)
	if len(runes) > maxLength {
		return fmt.Sprintf("%s...", string(runes[0:maxLength-3])), nil
	}
	return s, nil
}

// redraw redraws the screen.
func redraw() {
	//x, y := termbox.Size()

	drawBuildState(build{"test really really really really really really really really really really really really long", BuildStateFailed, false, ""}, rect{1, 1, 30, 3})
	drawBuildState(build{"test", BuildStatePassed, true, "Dave"}, rect{1, 5, 30, 3})
	drawBuildState(build{"test", BuildStateAcknowledged, false, "Sam"}, rect{1, 9, 30, 3})

	termbox.Flush()
}

func Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
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
