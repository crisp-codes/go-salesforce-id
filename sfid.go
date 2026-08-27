package sfid

import "errors"

var ErrInvalidID = errors.New("invalid id")

type ID string

func Parse(s string) (ID, error) {
	// TODO
	return ID(s), nil
}

func IsValid(id string) bool {
	// TODO
	return false
}

func To18Char(id string) (string, error) {
	// TODO
	return "", nil
}

func To15Char(id string) (string, error) {
	// TODO
	return "", nil
}

func (id ID) EntityPrefix() string {
	if len(id) < 3 {
		return ""
	}

	return string(id[:3])
}

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
