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
