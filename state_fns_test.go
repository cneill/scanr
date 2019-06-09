package scanr

import (
	"testing"
)

type stateFnTest struct {
	input string
	valid bool
}

var ipTests = []stateFnTest{
	{"127.0.0.1", true},
	{"999.999.999.999", false},
}

func TestIPScanr(t *testing.T) {
	for _, test := range ipTests {
		s := NewScanr(ScanIP)
		go s.Run(test.input)
	}
}
