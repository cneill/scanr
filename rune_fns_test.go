package scanr

import (
	"fmt"
	"strings"
	"testing"
)

type runeFnTest struct {
	input string
	valid bool
}

var hostnameTests = []runeFnTest{
	{"www.example.com", true},
	{"f49j0afj49jf40.com", true},
	{"this.is.a.long.subdomain.path.but.should.still.work.com", true},
	{"TECHNICALLY.THIS.WOULD.WORK.TOO.I.GUESS.com", true}, // TODO: but should it

	{"_ha24jfgik_.com", false},
	{"%20f", false},
	{"no spaces", false},

	{"111.111.111.111", true},        // this technically works too.. though it's not a hostname
	{"-invalid-but-valid.com", true}, // though the domain is not valid, it only contains valid hostname characters
}

var wordTests = []runeFnTest{
	{"this is a test", true},
	{"This Is A Test", true},
	{"THIS IS A TEST", true},
	{"ThIs Is A TeSt", true},

	{"Test 1234", false},
	{"This is a test.", false},
	{"This is a test !", false},
}

func TestHostnames(t *testing.T) {
	for _, test := range hostnameTests {
		var parts = strings.Split(test.input, ".")
		if err := testPartsValid(test, IsHostnameChar, parts); err != nil {
			t.Errorf("error: %v", err)
		}
	}
}

func TestWords(t *testing.T) {
	for _, test := range wordTests {
		var parts = strings.Split(test.input, " ")
		if err := testPartsValid(test, IsAlpha, parts); err != nil {
			t.Errorf("error: %v", err)
		}
	}
}

func testPartsValid(test runeFnTest, fn RuneFn, input []string) error {
	var anyInvalid = false
	for _, str := range input {
		if !testStringValid(str, fn) {
			anyInvalid = true
			break
		}
	}
	if test.valid && anyInvalid {
		return fmt.Errorf("valid test detected as invalid: %+v", test)
	} else if !test.valid && !anyInvalid {
		return fmt.Errorf("invalid test detected as valid: %+v", test)
	}
	return nil
}

func testStringValid(input string, fn RuneFn) bool {
	for _, c := range input {
		if !fn(rune(c)) {
			return false
		}
	}
	return true
}
