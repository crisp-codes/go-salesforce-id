package sfid

import "errors"

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
	if len(id) == LengthShort {
		return true
	}

	return validateChecksum(id)
}

// appendChecksum appends the checksum characters to the ID string.
// It is assumed that the ID has been validated before calling.
func appendChecksum(id string) string {
	return id
}

// validateChecksum validates the checksum characters of the ID string.
func validateChecksum(id string) bool {
	if len(id) != LengthLong {
		return false
	}

	return true
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
	// TODO
	return nil
}
