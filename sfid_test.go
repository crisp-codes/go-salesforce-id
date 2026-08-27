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
