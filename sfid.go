package sfid

import (
	"encoding/json"
	"errors"
)

const (
	LengthShort        = 15
	LengthLong         = 18
	LengthEntityPrefix = 3
)

const (
	checksumLookup = "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345"
)

var ErrInvalidID = errors.New("invalid id")

// ID represents a valid Salesforce record ID.
// The string representation of this type will always be the 18 character format.
type ID string

// Parse parses the provided string as an ID.
// Returns an error if the string is not in a valid Salesforce ID format.
func Parse(s string) (ID, error) {
	if !IsValid(s) {
		return "", ErrInvalidID
	}

	if len(s) == LengthShort {
		s, err := To18Char(s)
		if err != nil {
			return "", err
		}

		return ID(s), nil
	}

	return ID(s), nil
}

// IsValid checks that the provided string is in a valid Salesforce ID format.
// Both 15 char and 18 char formats are considered valid.
func IsValid(id string) bool {
	if len(id) != LengthShort && len(id) != LengthLong {
		return false
	}

	var block [3]int
	for i := range LengthShort {
		c := id[i]
		// SF IDs must be alphanumeric
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}

		if c >= 'A' && c <= 'Z' {
			block[i/5] += 1 << (i % 5)
		}
	}

	if len(id) == LengthShort {
		return true
	}

	for i, b := range block {
		if id[15+i] != checksumLookup[b] {
			return false
		}
	}

	return true
}

// appendChecksum appends the checksum characters to the ID string.
// It is assumed that the ID has been validated before calling.
func appendChecksum(id string) string {
	var block [3]int
	for i := range LengthShort {
		c := id[i]
		if c >= 'A' && c <= 'Z' {
			block[i/5] += 1 << (i % 5)
		}
	}

	checksum := ""
	for _, b := range block {
		checksum += string(checksumLookup[b])
	}

	return id + checksum
}

// To18Char validates and then returns the ID in 18 character string format
func To18Char(id string) (string, error) {
	if !IsValid(id) {
		return "", ErrInvalidID
	}
	if len(id) == LengthLong {
		return id, nil
	}

	return appendChecksum(id), nil
}

// To15Char validates and then returns the ID in 15 character string format
func To15Char(id string) (string, error) {
	if !IsValid(id) {
		return "", ErrInvalidID
	}

	return string(id[:LengthShort]), nil
}

// IsValid returns true if the ID is valid.
func (id ID) IsValid() bool {
	return IsValid(string(id))
}

// EntityPrefix returns the entity prefix of the id (e.g. "003" for contacts)
func (id ID) EntityPrefix() string {
	if len(id) < LengthEntityPrefix {
		return ""
	}

	return string(id[:LengthEntityPrefix])
}

// To15Char returns the id as a 15 character string
func (id ID) To15Char() string {
	if len(id) < LengthShort {
		return ""
	}

	return string(id[:LengthShort])
}

// Equal return true if the IDs are equal.
// It compares the 15 char case sensitive string representation of the IDs.
func (id ID) Equal(other ID) bool {
	return id.To15Char() == other.To15Char()
}

// IsZero returns true if the ID is empty.
// An invalid ID is not considered zero.
func (id ID) IsZero() bool {
	return string(id) == ""
}

// String returns the string literal of the ID
func (id ID) String() string {
	return string(id)
}

func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(id) + `"`), nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsed, err := Parse(s)
	if err != nil {
		return err
	}

	*id = parsed

	return nil
}
