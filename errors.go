package personnummer

// FormatError is returned by Parse and Personnummer.Parse if the supplied
// argument does not meet the basic format requirements.
type FormatError string

func (e FormatError) Error() string {
	return "invalid format: " + string(e)
}
