package sfid

import "errors"

var ErrInvalidID = errors.New("invalid id")

// ID represents a valid Salesforce record ID.
type ID string

// Parse parses the provided string as an ID.
// Returns an error if the string is not in a valid Salesforce ID format.
func Parse(s string) (ID, error) {
	// TODO
	return ID(s), nil
}

// IsValid checks that the provided string is in a valid Salesforce ID format.
// Both 15 char and 18 char formats are considered valid.
func IsValid(id string) bool {
	// TODO
	return false
}

// To18Char validates and then returns the ID in 18 character string format
func To18Char(id string) (string, error) {
	// TODO
	return "", nil
}

// To15Char validates and then returns the ID in 15 character string format
func To15Char(id string) (string, error) {
	// TODO
	return "", nil
}

// EntityPrefix returns the entity prefix of the id (e.g. "003" for contacts)
func (id ID) EntityPrefix() string {
	if len(id) < 3 {
		return ""
	}

	return string(id[:3])
}

// To15Char returns the id as a 15 character string
func (id ID) To15Char() string {
	if len(id) < 15 {
		return ""
	}

	return string(id[:15])
}

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
