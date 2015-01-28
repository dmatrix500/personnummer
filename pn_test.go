package personnummer

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

type pt struct {
	formats  []string
	date     time.Time
	nummer   int
	checksum int
	valid    bool
	female   bool
}

var testCases = []pt{
	// copy & paste for new cases:
	// pt{[]string{"", "", ""}, time.Date(, , , 0, 0, 0, 0, time.UTC), , , true, true}, //

	// VALID personnummer male
	pt{[]string{"196010052030", "6010052030", "19601005-2030", "601005-2030"}, time.Date(1960, 10, 5, 0, 0, 0, 0, time.UTC), 203, 0, true, false}, // 00
	pt{[]string{"197803164651", "7803164651", "19780316-4651", "780316-4651"}, time.Date(1978, 3, 16, 0, 0, 0, 0, time.UTC), 465, 1, true, false}, // 01
	pt{[]string{"195602241738", "5602241738", "19560224-1738", "560224-1738"}, time.Date(1956, 2, 24, 0, 0, 0, 0, time.UTC), 173, 8, true, false}, // 02
	pt{[]string{"199003175818", "9003175818", "19900317-5818", "900317-5818"}, time.Date(1990, 3, 17, 0, 0, 0, 0, time.UTC), 581, 8, true, false}, // 03
	pt{[]string{"198803255853", "8803255853", "19880325-5853", "880325-5853"}, time.Date(1988, 3, 25, 0, 0, 0, 0, time.UTC), 585, 3, true, false}, // 04
	pt{[]string{"196005135212", "6005135212", "19600513-5212", "600513-5212"}, time.Date(1960, 5, 13, 0, 0, 0, 0, time.UTC), 521, 2, true, false}, // 05
	pt{[]string{"198006102076", "8006102076", "19800610-2076", "800610-2076"}, time.Date(1980, 6, 10, 0, 0, 0, 0, time.UTC), 207, 6, true, false}, // 06
	pt{[]string{"197709136357", "7709136357", "19770913-6357", "770913-6357"}, time.Date(1977, 9, 13, 0, 0, 0, 0, time.UTC), 635, 7, true, false}, // 07
	pt{[]string{"198809098752", "8809098752", "19880909-8752", "880909-8752"}, time.Date(1988, 9, 9, 0, 0, 0, 0, time.UTC), 875, 2, true, false},  // 08
	pt{[]string{"197506200091", "7506200091", "19750620-0091", "750620-0091"}, time.Date(1975, 6, 20, 0, 0, 0, 0, time.UTC), 9, 1, true, false},   // 08

	// VALID personnummer female
	pt{[]string{"198011166009", "8011166009", "19801116-6009", "801116-6009"}, time.Date(1980, 11, 16, 0, 0, 0, 0, time.UTC), 600, 9, true, true}, // 09
	pt{[]string{"196204217142", "6204217142", "19620421-7142", "620421-7142"}, time.Date(1962, 4, 21, 0, 0, 0, 0, time.UTC), 714, 2, true, true},  // 10
	pt{[]string{"198212221546", "8212221546", "19821222-1546", "821222-1546"}, time.Date(1982, 12, 22, 0, 0, 0, 0, time.UTC), 154, 6, true, true}, // 11
	pt{[]string{"195409191847", "5409191847", "19540919-1847", "540919-1847"}, time.Date(1954, 9, 19, 0, 0, 0, 0, time.UTC), 184, 7, true, true},  // 12
	pt{[]string{"198701236088", "8701236088", "19870123-6088", "870123-6088"}, time.Date(1987, 1, 23, 0, 0, 0, 0, time.UTC), 608, 8, true, true},  // 13
	pt{[]string{"196811309860", "6811309860", "19681130-9860", "681130-9860"}, time.Date(1968, 11, 30, 0, 0, 0, 0, time.UTC), 986, 0, true, true}, // 14
	pt{[]string{"199112166245", "9112166245", "19911216-6245", "911216-6245"}, time.Date(1991, 12, 16, 0, 0, 0, 0, time.UTC), 624, 5, true, true}, // 15
	pt{[]string{"198005166767", "8005166767", "19800516-6767", "800516-6767"}, time.Date(1980, 5, 16, 0, 0, 0, 0, time.UTC), 676, 7, true, true},  // 16
	pt{[]string{"196809244285", "6809244285", "19680924-4285", "680924-4285"}, time.Date(1968, 9, 24, 0, 0, 0, 0, time.UTC), 428, 5, true, true},  // 17
	pt{[]string{"195311027089", "5311027089", "19531102-7089", "531102-7089"}, time.Date(1953, 11, 02, 0, 0, 0, 0, time.UTC), 708, 9, true, true}, // 18
	pt{[]string{"195311258387", "5311258387", "19531125-8387", "531125-8387"}, time.Date(1953, 11, 25, 0, 0, 0, 0, time.UTC), 838, 7, true, true}, // 19

	// INVALID cases
	// incorrect checksum
	//pt{[]string{"200603012224", "0603012224", "20060301-2224", "0603012224"}, time.Date(2006, 03, 01, 0, 0, 0, 0, time.UTC), 222, 0, false, true}, //
}

func testParseIdx(t *testing.T, pts []pt, format int) {
	for i, pt := range pts {
		pn, err := Parse(pt.formats[format])

		if err != nil && pt.valid {
			t.Errorf("Failed to parse valid: %d: %v %s", i, pt, err)
			continue
		}

		if err == nil && !pt.valid {
			t.Errorf("Successfully parsed invalid %d: %v", i, pt)
			continue
		}

		if !Compare(t, pn, pt) {
			t.Errorf("Wrongly parsed valid %d: %+v %+v", i, pt, pn)
		}
	}
}

func Compare(t *testing.T, pn Personnummer, p pt) bool {
	if pn.Nummer() != p.nummer {
		t.Errorf("nummer not matching: should be %3d is %3d", p.nummer, pn.Nummer())
		return false
	}
	if pn.(*pnImpl).checksum != p.checksum {
		t.Errorf("checksum not matching: should be %1d is %1d", p.checksum, pn.(*pnImpl).checksum)
		return false
	}
	if pn.IsFemale() != p.female {
		t.Errorf("sec not matching: should be  %t is %t", p.female, pn.IsFemale())
		return false
	}
	if pn.BirthDate() != p.date {
		t.Errorf("date not matching: should be %s is %s", p.date, pn.BirthDate())
		return false
	}
	if pn.Valid() != p.valid {
		t.Logf("Checksum should be %d is %d", p.checksum, calcChecksum(pn.BirthDate(), pn.Nummer()))
		t.Errorf("validity not matching: should be %t is %t", p.valid, pn.Valid())
		return false
	}

	return true
}

func TestCalcSum(t *testing.T) {
	for i, pt := range testCases {
		cS := calcChecksum(pt.date, pt.nummer)
		if cS != pt.checksum {
			t.Fatalf("Failed to calculate checksum %d: should be %d but is %d", i, pt.checksum, cS)
		}
	}
}

func TestPNParseLong(t *testing.T) {
	testParseIdx(t, testCases, 0)
}

func TestPNParseShort(t *testing.T) {
	testParseIdx(t, testCases, 1)
}

func TestPNParseDashedLong(t *testing.T) {
	testParseIdx(t, testCases, 2)
}

func TestPNParseDashedShort(t *testing.T) {
	testParseIdx(t, testCases, 3)
}

func TestNewPN(t *testing.T) {
	for i, pt := range testCases {
		pn := NewPersonnummer(pt.date, pt.nummer)
		if !Compare(t, pn, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pn)
		}
	}
}

func TestString(t *testing.T) {
	for i, pt := range testCases {
		pn := NewPersonnummer(pt.date, pt.nummer)
		if !Compare(t, pn, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pn)
		}
		if pn.String() != pt.formats[2] {
			t.Fatalf("String returned wrong: should be %s is %s", pn.String(), pt.formats[2])
		}
	}

}

func TestJSON(t *testing.T) {
	for i, pt := range testCases {
		pn := NewPersonnummer(pt.date, pt.nummer)
		if !Compare(t, pn, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pn)
		}
		b := bytes.NewBuffer(nil)
		en := json.NewEncoder(b)
		if err := en.Encode(pn); err != nil {
			t.Fatalf("Failed to encode %d to JSON %s (content: %s)", i, err, string(b.Bytes()))
		}
		dec := json.NewDecoder(b)

		pnDec := NewPersonnummer(time.Now(), 123)
		if err := dec.Decode(&pnDec); err != nil {
			t.Fatalf("Failed to decode %d from JSON %s (content: %s)", i, err, string(b.Bytes()))
		}

		if !Compare(t, pnDec, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pnDec)
		}
	}
}

func TestTEXTDeEn(t *testing.T) {
	for i, pt := range testCases {
		pn := NewPersonnummer(pt.date, pt.nummer)
		if !Compare(t, pn, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pn)
		}
		bt, err := pn.MarshalText()
		if err != nil {
			t.Fatalf("Failed to encode %d to Text %s (content: %s)", i, err, string(bt))
		}

		pnDec := NewPersonnummer(time.Now(), 123)
		if err := pnDec.UnmarshalText(bt); err != nil {
			t.Fatalf("Failed to decode %d from Text %s (content: %s)", i, err, string(bt))
		}

		if !Compare(t, pnDec, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pnDec)
		}
	}
}

func TestSQL(t *testing.T) {
	for i, pt := range testCases {
		pn := NewPersonnummer(pt.date, pt.nummer)
		if !Compare(t, pn, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pn)
		}
		ie, err := pn.Value()
		if err != nil {
			t.Fatalf("Failed to encode %d to SQL %s (content: %v)", i, err, ie)
		}

		pnDec := NewPersonnummer(time.Now(), 123)
		if err := pnDec.Scan(ie); err != nil {
			t.Fatalf("Failed to decode %d from Text %s (content: %v)", i, err, ie)
		}

		if !Compare(t, pnDec, pt) {
			t.Fatalf("NewPersonnummer failed %d: %v", i, pnDec)
		}
	}
}
