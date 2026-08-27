package sfid_test

import (
	"errors"
	"testing"

	sfid "github.com/crisp-codes/go-salesforce-id"
)

func FuzzParse(f *testing.F) {
	for _, test := range parseTests {
		f.Add(test.input)
	}
	for _, test := range isValidTests {
		f.Add(test.input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		id, err := sfid.Parse(s)

		if sfid.IsValid(s) {
			if err != nil {
				t.Fatalf("Parse(%q) returned unexpected error: %v", s, err)
			}
			if len(id) != sfid.LengthLong {
				t.Fatalf("Parse(%q) = %q, want length %d", s, id, sfid.LengthLong)
			}
			if !id.IsValid() {
				t.Fatalf("Parse(%q) = %q, which is not itself valid", s, id)
			}

			shortID, err := sfid.To15Char(s)
			if err != nil {
				t.Fatalf("To15Char(%q) returned unexpected error: %v", s, err)
			}
			if shortID != id.To15Char() {
				t.Fatalf("To15Char(%q) = %q, want %q", s, shortID, id.To15Char())
			}

			reparse, err := sfid.Parse(shortID)
			if err != nil {
				t.Fatalf("Parse(%q) returned unexpected error: %v", shortID, err)
			}
			if reparse != id {
				t.Fatalf("Parse(%q) = %q, want %q", shortID, reparse, id)
			}
		} else {
			if !errors.Is(err, sfid.ErrInvalidID) {
				t.Fatalf("Parse(%q) expected ErrInvalidID, got: %v", s, err)
			}
			if id != "" {
				t.Fatalf("Parse(%q) expected empty ID on error, got %q", s, id)
			}
		}
	})
}
