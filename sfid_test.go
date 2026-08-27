package sfid_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"crisp.codes/sfid"
)

var (
	_ fmt.Stringer     = sfid.ID("")
	_ json.Marshaler   = sfid.ID("")
	_ json.Unmarshaler = (*sfid.ID)(nil)
)

var parseTests = []struct {
	name     string
	input    string
	isValid  bool
	expected sfid.ID
}{
	{
		name:     "18 char valid",
		input:    `012000000000000AAA`,
		isValid:  true,
		expected: sfid.ID("012000000000000AAA"),
	},
	{
		name:    "18 char invalid",
		input:   `012000000000000ZZZ`,
		isValid: false,
	},
	{
		name:     "15 char valid",
		input:    `012000000000000`,
		isValid:  true,
		expected: sfid.ID("012000000000000AAA"),
	},
	{
		name:    "empty",
		input:   ``,
		isValid: false,
	},
}

func Test_Parse(t *testing.T) {
	for _, test := range parseTests {
		t.Run(test.name, func(t *testing.T) {
			id, err := sfid.Parse(test.input)
			if test.isValid {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if id != test.expected {
					t.Errorf("unexpected ID: %v", id)
				}
			} else {
				if !errors.Is(err, sfid.ErrInvalidID) {
					t.Errorf("expected invalid ID error, got: %v", err)
				}
			}
		})
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	for _, test := range parseTests {
		t.Run(test.name, func(t *testing.T) {
			var id sfid.ID
			err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, test.input)), &id)
			if test.isValid {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if id != test.expected {
					t.Errorf("unexpected ID: %v", id)
				}
			} else {
				if !errors.Is(err, sfid.ErrInvalidID) {
					t.Errorf("expected invalid ID error, got: %v", err)
				}
			}
		})
	}
}

var isValidTests = []struct {
	name  string
	input string
	valid bool
}{
	{name: "18 char valid, all digits", input: "012000000000000AAA", valid: true},
	{name: "18 char valid, all uppercase", input: "AAAAAAAAAAAAAAA555", valid: true},
	{name: "18 char valid, mixed case", input: "a0Ab00000FooBaREQU", valid: true},
	{name: "18 char valid, contact prefix", input: "003AaBbCcDdEeF1IVK", valid: true},
	{name: "18 char invalid checksum", input: "012000000000000ZZZ", valid: false},
	{name: "15 char valid, all digits", input: "012000000000000", valid: true},
	{name: "15 char valid, mixed case", input: "a0Ab00000FooBaR", valid: true},
	{name: "too short", input: "01200000000000", valid: false},
	{name: "too long", input: "012000000000000AAAA", valid: false},
	{name: "17 chars", input: "012000000000000AA", valid: false},
	{name: "empty", input: "", valid: false},
	{name: "non-alphanumeric in body", input: "01200000000000!", valid: false},
	{name: "non-alphanumeric in checksum", input: "012000000000000AA!", valid: false},
	{name: "whitespace", input: "               ", valid: false},
}

func Test_IsValid(t *testing.T) {
	for _, test := range isValidTests {
		t.Run(test.name, func(t *testing.T) {
			if got := sfid.IsValid(test.input); got != test.valid {
				t.Errorf("IsValid(%q) = %v, want %v", test.input, got, test.valid)
			}
		})
	}
}

func TestID_IsValid(t *testing.T) {
	for _, test := range isValidTests {
		t.Run(test.name, func(t *testing.T) {
			id := sfid.ID(test.input)
			if got := id.IsValid(); got != test.valid {
				t.Errorf("ID(%q).IsValid() = %v, want %v", test.input, got, test.valid)
			}
		})
	}
}

var to18CharTests = []struct {
	name     string
	input    string
	err      bool
	expected string
}{
	{name: "15 char, all digits", input: "012000000000000", expected: "012000000000000AAA"},
	{name: "15 char, all uppercase", input: "AAAAAAAAAAAAAAA", expected: "AAAAAAAAAAAAAAA555"},
	{name: "15 char, mixed case", input: "a0Ab00000FooBaR", expected: "a0Ab00000FooBaREQU"},
	{name: "18 char passthrough", input: "012000000000000AAA", expected: "012000000000000AAA"},
	{name: "invalid", input: "not-a-valid-id!", err: true},
	{name: "empty", input: "", err: true},
}

func Test_To18Char(t *testing.T) {
	for _, test := range to18CharTests {
		t.Run(test.name, func(t *testing.T) {
			got, err := sfid.To18Char(test.input)
			if test.err {
				if !errors.Is(err, sfid.ErrInvalidID) {
					t.Errorf("expected invalid ID error, got: %v", err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if got != test.expected {
				t.Errorf("To18Char(%q) = %q, want %q", test.input, got, test.expected)
			}
		})
	}
}

var to15CharTests = []struct {
	name     string
	input    string
	err      bool
	expected string
}{
	{name: "15 char passthrough", input: "012000000000000", expected: "012000000000000"},
	{name: "18 char truncated", input: "012000000000000AAA", expected: "012000000000000"},
	{name: "18 char, mixed case", input: "a0Ab00000FooBaREQU", expected: "a0Ab00000FooBaR"},
	{name: "invalid checksum", input: "012000000000000ZZZ", err: true},
	{name: "invalid length", input: "012", err: true},
	{name: "empty", input: "", err: true},
}

func Test_To15Char(t *testing.T) {
	for _, test := range to15CharTests {
		t.Run(test.name, func(t *testing.T) {
			got, err := sfid.To15Char(test.input)
			if test.err {
				if !errors.Is(err, sfid.ErrInvalidID) {
					t.Errorf("expected invalid ID error, got: %v", err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if got != test.expected {
				t.Errorf("To15Char(%q) = %q, want %q", test.input, got, test.expected)
			}
		})
	}
}

func TestID_To15Char(t *testing.T) {
	tests := []struct {
		name     string
		input    sfid.ID
		expected string
	}{
		{name: "18 char truncated", input: sfid.ID("012000000000000AAA"), expected: "012000000000000"},
		{name: "15 char unchanged", input: sfid.ID("012000000000000"), expected: "012000000000000"},
		{name: "too short returns empty", input: sfid.ID("012"), expected: ""},
		{name: "empty returns empty", input: sfid.ID(""), expected: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.input.To15Char(); got != test.expected {
				t.Errorf("To15Char() = %q, want %q", got, test.expected)
			}
		})
	}
}

func TestID_EntityPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    sfid.ID
		expected string
	}{
		{name: "contact prefix", input: sfid.ID("003AaBbCcDdEeF1IVK"), expected: "003"},
		{name: "15 char id", input: sfid.ID("012000000000000"), expected: "012"},
		{name: "exactly 3 chars", input: sfid.ID("003"), expected: "003"},
		{name: "shorter than prefix", input: sfid.ID("01"), expected: ""},
		{name: "empty", input: sfid.ID(""), expected: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.input.EntityPrefix(); got != test.expected {
				t.Errorf("EntityPrefix() = %q, want %q", got, test.expected)
			}
		})
	}
}

func TestID_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a        sfid.ID
		b        sfid.ID
		expected bool
	}{
		{
			name:     "same 18 char id",
			a:        sfid.ID("003AaBbCcDdEeF1IVK"),
			b:        sfid.ID("003AaBbCcDdEeF1IVK"),
			expected: true,
		},
		{
			name:     "15 char equal to its own 18 char form",
			a:        sfid.ID("012000000000000"),
			b:        sfid.ID("012000000000000AAA"),
			expected: true,
		},
		{
			name:     "different ids",
			a:        sfid.ID("003AaBbCcDdEeF1IVK"),
			b:        sfid.ID("012000000000000AAA"),
			expected: false,
		},
		{
			name:     "differs only in checksum suffix is still equal",
			a:        sfid.ID("012000000000000AAA"),
			b:        sfid.ID("012000000000000ZZZ"),
			expected: true,
		},
		{
			name:     "case sensitive on 15 char body",
			a:        sfid.ID("a0Ab00000FooBaR"),
			b:        sfid.ID("A0AB00000FOOBAR"),
			expected: false,
		},
		{
			name:     "both empty",
			a:        sfid.ID(""),
			b:        sfid.ID(""),
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.a.Equal(test.b); got != test.expected {
				t.Errorf("%q.Equal(%q) = %v, want %v", test.a, test.b, got, test.expected)
			}
		})
	}
}

func TestID_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		input    sfid.ID
		expected bool
	}{
		{name: "empty is zero", input: sfid.ID(""), expected: true},
		{name: "valid id is not zero", input: sfid.ID("012000000000000AAA"), expected: false},
		{name: "invalid non-empty string is not zero", input: sfid.ID("not-valid"), expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.input.IsZero(); got != test.expected {
				t.Errorf("IsZero() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestID_String(t *testing.T) {
	tests := []struct {
		name     string
		input    sfid.ID
		expected string
	}{
		{name: "18 char id", input: sfid.ID("012000000000000AAA"), expected: "012000000000000AAA"},
		{name: "empty id", input: sfid.ID(""), expected: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.input.String(); got != test.expected {
				t.Errorf("String() = %q, want %q", got, test.expected)
			}
		})
	}
}

func TestID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    sfid.ID
		expected string
	}{
		{name: "18 char id", input: sfid.ID("012000000000000AAA"), expected: `"012000000000000AAA"`},
		{name: "empty id", input: sfid.ID(""), expected: `""`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if string(b) != test.expected {
				t.Errorf("MarshalJSON() = %s, want %s", b, test.expected)
			}
		})
	}
}

func TestID_MarshalJSON_struct(t *testing.T) {
	type record struct {
		ID sfid.ID `json:"id"`
	}

	in := record{ID: sfid.ID("012000000000000AAA")}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out record
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !in.ID.Equal(out.ID) {
		t.Errorf("round trip mismatch: got %v, want %v", out.ID, in.ID)
	}
}

func TestID_UnmarshalJSON_invalidJSON(t *testing.T) {
	var id sfid.ID
	err := json.Unmarshal([]byte(`123`), &id)
	if err == nil {
		t.Error("expected error unmarshaling non-string JSON value")
	}
}
