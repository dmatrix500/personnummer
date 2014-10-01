package personnummer

import (
	"time"
)

type pnImpl struct {
	birth    time.Time
	nummer   int
	checksum int
}

// NewPersonnummer returns a new personnummer from the given birthdate and
// running number with the checksum initialized accordingly.
func NewPersonnummer(birthdate time.Time, runningNummer int) Personnummer {
	p := &pnImpl{
		birth:  birthdate,
		nummer: runningNummer,
	}
	p.checksum = calcChecksum(p.birth, p.nummer)
	return p
}

// Valid returns true if the personnummer is valid.
func (p pnImpl) Valid() bool {
	return calcChecksum(p.birth, p.nummer) == p.checksum
}

// IsFemale returns true if personnummber describes a female person, false
// if it describes a male person
func (p pnImpl) IsFemale() bool {
	return p.nummer%2 == 0
}

// BirthDate returns the birth date of the person
func (p pnImpl) BirthDate() time.Time {
	return p.birth
}

func (p pnImpl) Nummer() int {
	return p.nummer
}
