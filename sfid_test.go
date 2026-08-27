package sfid_test

import (
	"encoding/json"
	"fmt"

	"crisp.codes/sfid"
)

var (
	_ fmt.Stringer     = sfid.ID("")
	_ json.Marshaler   = sfid.ID("")
	_ json.Unmarshaler = (*sfid.ID)(nil)
)
