package personnummer

import (
	"time"
)

type pnImpl struct {
	Birth           time.Time
	RunningNumber   int
	Checksum        int
	NumberSeparator string
}

// NewPersonnummer returns a new personnummer from the given birthdate and
// running number with the checksum initialized accordingly. Use Parse(string)
// to parse a comlete personnummer.
func NewPersonnummer(birthdate time.Time, runningNummer int) Personnummer {
	p := &pnImpl{
		Birth:         birthdate,
		RunningNumber: runningNummer,
	}
	p.Checksum = calcChecksum(p.Birth, p.RunningNumber)
	return p
}

func NewEmptyPN() Personnummer {
	return &pnImpl{}
}

// Valid returns true if the personnummer is valid.
func (p pnImpl) Valid() bool {
	return calcChecksum(p.Birth, p.RunningNumber) == p.Checksum
}

// IsFemale returns true if the person is female and
// false if the person is male.
func (p pnImpl) IsFemale() bool {
	return p.RunningNumber%2 == 0
}

// BirthDate returns the person's birth date.
func (p pnImpl) BirthDate() time.Time {
	return p.Birth
}

// Age returns the person's Age in years.
func (p pnImpl) Age() int {
	return int(time.Since(p.Birth).Hours() / 8765.81)
}

func (p pnImpl) Nummer() int {
	return p.RunningNumber
}

// Separator specifies which character separates the parts of the ssn
// during unmarshalling. Defaults to "".
func (p *pnImpl) Separator(s string) {
	p.NumberSeparator = s
}
