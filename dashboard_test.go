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
