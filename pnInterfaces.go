package personnummer

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Personnummer implements the swedish personnummer.
type Personnummer interface {
	// Parses the argument as a personnummer.
	Parse(string) error

	// Valid returns true if the personnummer is valid.
	Valid() bool

	// Returns true if personnummer specifies a female person.
	// False indicates a male person.
	IsFemale() bool

	// Returns the birthdate contained in the person nummer.
	BirthDate() time.Time

	// Returns the running number contained in the person nummer.
	Nummer() int

	// fmt.Stringer
	String() string

	// json.(Un)marshaller
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error

	// encoding.Text(Un)marshaller
	MarshalText() ([]byte, error)
	UnmarshalText([]byte) error

	// driver.Valuer, ql.Scanner
	Value() (driver.Value, error)
	Scan(interface{}) error
}

// pnImpl implements fmt.Stringer
func (p pnImpl) String() string {
	return fmt.Sprintf("%s-%d%d", p.birth.Format("20060102"), p.nummer, p.checksum)
}

// Value implements sql/driver.Valuer.
func (p pnImpl) Value() (driver.Value, error) {
	return p.String(), nil
}

// pnImpl implements sql.Scanner
func (p *pnImpl) Scan(src interface{}) error {
	var strVal string
	switch src.(type) {
	case string:
		strVal = src.(string)
	case []byte:
		strVal = string(src.([]byte))
	default:
		return fmt.Errorf("pnImpl.Scan: cannot convert %T to string.", src)
	}

	return p.Parse(strVal)
}

// MarshalText implements encoding.TextMarshaler.
func (p pnImpl) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *pnImpl) UnmarshalText(text []byte) error {
	return p.Parse(string(text))
}

// MarshalJSON implements json.Marshaler.
func (p pnImpl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p.String() + "\""), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *pnImpl) UnmarshalJSON(text []byte) error {
	return p.Parse(string(text[1 : len(text)-1]))
}
