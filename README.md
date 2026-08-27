# go-salesforce-id

A simple Go module providing typing and validation for Salesforce record IDs.

Handles both the 15 character (case-sensitive) and 18 character (case-insensitive) ID formats.

## Install

```sh
go get github.com/crisp-codes/go-salesforce-id
```

## Usage

```go
import "github.com/crisp-codes/go-salesforce-id"

id, err := sfid.Parse("003000000000001")
if err != nil {
    // ... do stuff ...
}

id.String() // "003000000000001AAA" (always normalized to 18 char)
id.To15Char() // "003000000000001"
id.EntityPrefix() // "003"

sfid.IsValid("003000000000001AAA") // true

a, _ := sfid.Parse("003000000000001")
b, _ := sfid.Parse("003000000000001AAA")
a.Equal(b) // true
```

`ID` also implements `json.Marshaler` / `json.Unmarshaler`, so it can be used directly as a struct field within API calls:

```go
type Contact struct {
    ID sfid.ID `json:"Id"`
}
```

## License

Apache 2.0, see [LICENSE](LICENSE).
