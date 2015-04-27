package monitrondashboard

// Basic dashboard functionality

import (
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
)

const OrangeColour int = 167

const textPadding int = 1

const buildingMessage string = "Building"

// buildState is an int type defining the states a build can be in.
type buildState int

const (
	BuildStateFailed buildState = iota
	BuildStateAcknowledged
	BuildStatePassed
	BuildStateUnknown
)

// BgColour returns the termbox Attribute for the background colour
// of a build in this state.
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

// BgColour returns the termbox Attribute for the text colour
// of a build in this state.
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

// build is a struct covering a build and its current state.
type build struct {
	name         string
	buildState   buildState
	building     bool
	acknowledger string
}

// rect is a simple struct giving a bounding rectangle for a widget
type rect struct {
	point
	size
}

func (r rect) String() string {
	return fmt.Sprintf("<Rect x:%d,y:%d | w:%d,h:%d>", r.x, r.y,
		r.w, r.h)
}

func NewRect(x, y, w, h int) rect {
	return rect{
		point{x, y},
		size{w, h},
	}
}

// point defines a single point on the screen
type point struct {
	x int
	y int
}

// size defines a width and height an object takes up
type size struct {
	w int
	h int
}

// drawBuildState draws a status box for an individual build within bounds
func drawBuildState(build build, bounds rect) {
	bgColour := build.buildState.BgColour()
	fgColour := build.buildState.FgColour()

	availableWidth := bounds.size.w - 2*textPadding
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

// Layout struct containing a slice of rectangles for each grid position.
type Layout struct {
	boxes []rect
}

// layoutGridForScreen returns a Layout detailing positioning for numberOfBoxes,
// taking into account the mininumBoxSize and padding, fitting onto screenSize.
func layoutGridForScreen(minimumBoxSize size, numberOfBoxes int, padding int, screenSize size) (Layout, error) {
	maximumNumberOfVerticalBoxes := (screenSize.h - padding) / (minimumBoxSize.h + padding)
	// integer division that always rounds up
	requiredNumberOfColumns := (numberOfBoxes + maximumNumberOfVerticalBoxes - 1) / maximumNumberOfVerticalBoxes
	columnWidth := (screenSize.w - padding) / requiredNumberOfColumns
	boxWidth := columnWidth - padding
	if boxWidth < minimumBoxSize.w {
		return Layout{}, errors.New("Screen is too small to fit the grid")
	}

	boxes := []rect{}
	i := 0
gridloop:
	for x := 0; x < requiredNumberOfColumns; x++ {
		for y := 0; y < maximumNumberOfVerticalBoxes; y++ {
			boxes = append(boxes, rect{
				point{
					padding + x*columnWidth,
					padding + y*(minimumBoxSize.h+padding),
				},
				size{boxWidth, minimumBoxSize.h},
			})
			i = i + 1
			if i == numberOfBoxes {
				break gridloop
			}
		}
	}

	return Layout{
		boxes: boxes,
	}, nil
}

// redraw redraws the screen.
func redraw() error {
	screenWidth, screenHeight := termbox.Size()
	layout, err := layoutGridForScreen(size{30, 3}, 15, 1, size{screenWidth, screenHeight})
	if err != nil {
		return err
	}

	for i := 0; i < 15; i++ {
		box := layout.boxes[i]
		drawBuildState(build{"test", BuildStateFailed, false, ""}, box)
	}

	termbox.Flush()
	return nil
}

func Run() {
	NewBuildFetcher()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
	if err := redraw(); err != nil {
		fmt.Println("Error: %s", err)
		return
	}

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
		case termbox.EventResize:
			termbox.Clear(0, 0)
			if err := redraw(); err != nil {
				fmt.Println("Error: %s", err)
				return
			}
		}
		if err := redraw(); err != nil {
			fmt.Println("Error: %s", err)
			return
		}
	}
}
