package monitrondashboard

// Tests for the Monitron dashboard

import "testing"

var ellipsizeTests = []struct {
	in          string
	maxLength   int
	out         string
	shouldError bool
}{
	{"short", 10, "short", false},
	{"exact length name", 17, "exact length name", false},
	{"long name", 6, "lon...", false},
	{"⌘ ⌘ ⌘", 5, "⌘ ⌘ ⌘", false},
	{"⌘ ⌘ ⌘ ⌘ ⌘ ⌘", 6, "⌘ ⌘...", false},
	{"length too short for elipsis", 2, "", true},
}

func TestEllipsize(t *testing.T) {
	for _, test := range ellipsizeTests {
		out, err := elipsize(test.in, test.maxLength)
		if test.shouldError {
			if err == nil {
				t.Errorf("ellipsize(%q, %d) should fail with Error",
					test.in, test.maxLength)
			}
		} else {
			if err != nil {
				t.Errorf("ellipsize(%q, %d) should not fail with Error: %s",
					test.in, test.maxLength, err)
			}
		}
		if out != test.out {
			t.Errorf("ellipsize(%q, %d) => %q, expected %q",
				test.in, test.maxLength, out, test.out)
		}
	}
}

func TestTextWriterPrintsLabelInCorrectPlace(t *testing.T) {
	writerToTest := createTextWriter("Label", point{1, 1})

	y := 0
	x := 0
	// Should not print on the first line
	for x := 0; x < 10; x++ {
		char := writerToTest(' ', point{x, y})
		if char != ' ' {
			t.Errorf("Label starting at (1, 1) should not print at: (%d, %d)",
				x, y)
		}
	}

	y = 1
	x = 0
	// Should not print on first column of second line
	char := writerToTest(' ', point{x, y})
	if char != ' ' {
		t.Errorf("Label starting at (1, 1) should not print at: (%d, %d)",
			x, y)
	}

	expectedRunes := []rune("Label")
	// Should print on the second line
	for x := 1; x < 6; x++ {
		char := writerToTest(' ', point{x, y})
		if char != expectedRunes[x-1] {
			t.Errorf("Label starting at (1, 1) should print %q at: (%d, %d), got: %q",
				expectedRunes[x-1], x, y, char)
		}
	}
}

func TestTextWriterShouldNotExceedItsBounds(t *testing.T) {

}

func TestLayoutGridCorrectlyLaysOutRectangles(t *testing.T) {
	numberOfBoxes := 9
	layout, err := layoutGridForScreen(size{300, 3}, numberOfBoxes, 1, size{904, 13})
	if err != nil {
		t.Fatalf("Unexpected error when calling layoutGridForScreen: %s", err)
	}

	if len(layout.boxes) != numberOfBoxes {
		t.Fatalf("layoutGridForScreen: should provide a box for every build, got %d boxes expected %d", len(layout.boxes), numberOfBoxes)
	}

	expectBoxDimensions := func(i int, box rect, expectedBox rect) {
		if box != expectedBox {
			t.Errorf("Box: %d was %s, expected: %s", i, box, expectedBox)
		}
	}
	expectBoxDimensions(0, layout.boxes[0], NewRect(1, 1, 300, 3))
	expectBoxDimensions(1, layout.boxes[1], NewRect(1, 5, 300, 3))
	expectBoxDimensions(2, layout.boxes[2], NewRect(1, 9, 300, 3))
	expectBoxDimensions(3, layout.boxes[3], NewRect(1+301, 1, 300, 3))
	expectBoxDimensions(4, layout.boxes[4], NewRect(1+301, 5, 300, 3))
	expectBoxDimensions(5, layout.boxes[5], NewRect(1+301, 9, 300, 3))
	expectBoxDimensions(6, layout.boxes[6], NewRect(1+301*2, 1, 300, 3))
	expectBoxDimensions(7, layout.boxes[7], NewRect(1+301*2, 5, 300, 3))
	expectBoxDimensions(8, layout.boxes[8], NewRect(1+301*2, 9, 300, 3))
}

func TestLayoutGridCorrectlyLaysOutRectanglesLeavingGapsForUnevenColumns(t *testing.T) {
	numberOfBoxes := 8
	layout, err := layoutGridForScreen(size{300, 3}, numberOfBoxes, 1, size{904, 13})
	if err != nil {
		t.Fatalf("Unexpected error when calling layoutGridForScreen: %s", err)
	}

	if len(layout.boxes) != numberOfBoxes {
		t.Fatalf("layoutGridForScreen: should provide a box for every build, got %d boxes expected %d", len(layout.boxes), numberOfBoxes)
	}

	expectBoxDimensions := func(i int, box rect, expectedBox rect) {
		if box != expectedBox {
			t.Errorf("Box: %d was %s, expected: %s", i, box, expectedBox)
		}
	}
	expectBoxDimensions(0, layout.boxes[0], NewRect(1, 1, 300, 3))
	expectBoxDimensions(1, layout.boxes[1], NewRect(1, 5, 300, 3))
	expectBoxDimensions(2, layout.boxes[2], NewRect(1, 9, 300, 3))
	expectBoxDimensions(3, layout.boxes[3], NewRect(1+301, 1, 300, 3))
	expectBoxDimensions(4, layout.boxes[4], NewRect(1+301, 5, 300, 3))
	expectBoxDimensions(5, layout.boxes[5], NewRect(1+301, 9, 300, 3))
	expectBoxDimensions(6, layout.boxes[6], NewRect(1+301*2, 1, 300, 3))
	expectBoxDimensions(7, layout.boxes[7], NewRect(1+301*2, 5, 300, 3))
}

func TestLayoutGridErrorsWhenNotEnoughSpace(t *testing.T) {
	numberOfBoxes := 9
	_, err := layoutGridForScreen(size{300, 3}, numberOfBoxes, 1, size{800, 13})
	if err == nil {
		t.Errorf("Expected error when calling layoutGridForScreen, when screen too small")
	}
}
