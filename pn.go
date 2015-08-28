package personnummer

import (
	"time"
)

type pnImpl struct {
	birth     time.Time
	nummer    int
	checksum  int
	separator string
}

// NewPersonnummer returns a new personnummer from the given birthdate and
// running number with the checksum initialized accordingly. Use Parse(string)
// to parse a comlete personnummer.
func NewPersonnummer(birthdate time.Time, runningNummer int) Personnummer {
	p := &pnImpl{
		birth:  birthdate,
		nummer: runningNummer,
	}
	p.checksum = calcChecksum(p.birth, p.nummer)
	return p
}

func NewEmptyPN() Personnummer {
	return &pnImpl{}
}

// Valid returns true if the personnummer is valid.
func (p pnImpl) Valid() bool {
	return calcChecksum(p.birth, p.nummer) == p.checksum
}

// IsFemale returns true if the person is female and
// false if the person is male.
func (p pnImpl) IsFemale() bool {
	return p.nummer%2 == 0
}

// BirthDate returns the person's birth date.
func (p pnImpl) BirthDate() time.Time {
	return p.birth
}

// Age returns the person's Age in years.
func (p pnImpl) Age() int {
	return int(time.Since(p.birth).Hours() / 8765.81)
}

func (p pnImpl) Nummer() int {
	return p.nummer
}

// Separator specifies which character separates the parts of the ssn
// during unmarshalling. Defaults to "".
func (p *pnImpl) Separator(s string) {
	p.separator = s
}
