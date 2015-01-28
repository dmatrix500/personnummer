package personnummer

import (
	"strconv"
	"strings"
	"time"
)

func calcChecksum(birth time.Time, nummer int) int {

	numSing := nummer % 10
	numCent := nummer / 100
	numDez := nummer/10 - 10*numCent
	// fmt.Printf("nummer %d: %d %d %d\n", nummer, numCent, numDez, numSing)

	var calcs = map[int][]int{
		1: []int{
			birth.Year() % 10,
			int(birth.Month()) % 10,
			birth.Day() % 10,
			numDez,
		},
		2: []int{
			(birth.Year() % 100) / 10,
			int(birth.Month()) / 10,
			birth.Day() / 10,
			numCent,
			numSing,
		},
	}
	//fmt.Printf("date: %s: %v %v\n", birth.Format("2006 01 02"), calcs[2], calcs[1])

	var sum int
	for multiplier, values := range calcs {
		for _, value := range values {
			prod := multiplier * value
			//fmt.Printf("mul %d: value %d prod %d\n", multiplier, value, prod)
			if prod > 9 {
				prod = prod/10 + prod%10
				//fmt.Printf("top %d: value %d prod %d\n", multiplier, value, prod)
			}
			sum += prod
		}
	}

	sum = 10 - sum%10
	if sum == 10 {
		sum = 0
	}
	return sum

}

func splitSign(pn string) (date string, nummer string, check string, err error) {
	parts := strings.Split(pn, "-")
	if len(parts) != 2 {
		err = FormatError("missing or to many '-' in personnummer")
		return
	}

	if len(parts[0]) != 8 && len(parts[0]) != 6 {
		err = FormatError("wrongly positioned sign in personnummer")
		return
	}

	date = parts[0]
	nummer = parts[1][:3]
	check = parts[1][3:]

	return
}

// parse parses the string pn as swedish personnummer. In case of an error
// the returned error is FormatError if the format does not fit or an error
// returned by the time.Parse or strconv.Atoi functions called during parsing.
func (p *pnImpl) Parse(pn string) error {
	var date, nummerString, checkString string
	var err error
	var shortDate bool

	switch len(pn) {
	case 12:
		// long year, no dash
		date, nummerString, checkString = pn[:8], pn[8:11], pn[11:]
	case 10:
		// short year, no dash
		date, nummerString, checkString = pn[:6], pn[6:9], pn[9:]
		shortDate = true
	case 11:
		// short year, with dash
		date, nummerString, checkString, err = splitSign(pn)
		shortDate = true
	case 13:
		// long year, with dash
		date, nummerString, checkString, err = splitSign(pn)
	default:
		return FormatError("wrong length, should be 10-13, but is " + strconv.Itoa(len(pn)))
	}

	if err != nil {
		return err
	}
	if !shortDate {
		if p.birth, err = time.Parse("20060102", date); err != nil {
			return err
		}
	} else {
		if p.birth, err = time.Parse("060102", date); err != nil {
			return err
		}
	}

	if p.birth.Year() > time.Now().Year() {
		p.birth = time.Date(p.birth.Year()-100, p.birth.Month(), p.birth.Day(), 0, 0, 0, 0, time.UTC)
	}

	if p.nummer, err = strconv.Atoi(nummerString); err != nil {
		return err
	}

	if p.checksum, err = strconv.Atoi(checkString); err != nil {
		return err
	}

	return nil
}

// Parse parses pn as a swedish personnummer. It is a forwarder for Personnummer.Parse.
// In case of an error the returned Personnumber will be nil.
func Parse(pn string) (Personnummer, error) {
	p := &pnImpl{}
	err := p.Parse(pn)
	if err != nil {
		return nil, err
	}
	return p, nil
}
