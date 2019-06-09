package scanr

import "fmt"

// StateFn is a function used to transition between states in the Scanr
type StateFn func(*Scanr) StateFn

// ScanSpace scans & emits a run of space characters
func ScanSpace(s *Scanr) StateFn {
	s.AcceptWhileRuneFn(IsSpace)
	s.Emit(ItemSpace)
	return s.homeState
}

// ScanNewline scans & emits a Newline of the variants "\n", "\r", and "\r\n"
func ScanNewline(s *Scanr) StateFn {
	n := s.Next()
	if !IsNewline(n) {
		s.Backup()
		return s.homeState
	} else if n == '\r' {
		s.Accept("\n")
	}

	s.Emit(ItemNewline)
	return s.homeState
}

func getIPOctet(s *Scanr) error {
	var leaderOne, leaderTwo, atMax2ndDigit bool

	// hundreds place, eg: 1XX.XXX... or 2XX.XXX
	if leaderOne = s.Accept("1"); !leaderOne {
		leaderTwo = s.Accept("2")
	}
	if leaderTwo {
		if !s.Accept("01234") {
			if s.Accept("5") {
				atMax2ndDigit = true
			}
		}

		if atMax2ndDigit {
			s.Accept("012345")
		} else {
			s.Accept(digits)
		}
	} else {
		s.Accept(digits)
		s.Accept(digits)
	}

	return nil
}

// ScanIP scans & emits an IP
func ScanIP(s *Scanr) StateFn {
	// TODO: better parsing of octets? or handle at the Parser level?
	// s.AcceptWhileRuneFn(IsNumber)
	for i := 0; i < 3; i++ {
		if err := getIPOctet(s); err != nil {
			s.Emit(ItemError)
		}
		if !s.Accept(".") {
			s.Emit(ItemError)
		}
	}
	if err := getIPOctet(s); err != nil {
		s.Emit(ItemError)
	}
	if p := s.Peek(); IsNumber(p) || p == '.' {
		// too many octets or too many digits
		s.Emit(ItemError)
	}
	s.Emit(ItemIP)
	return s.homeState
}

func getHostnamePart(s *Scanr) error {
	if n := s.Next(); !IsAlphaNum(n) {
		// domains must start with an alphanumeric character
		s.Emit(ItemError)
	}

	s.AcceptWhileRuneFn(IsHostnameChar)

	if s.PrevRune() == '-' {
		return fmt.Errorf("hostnames can't end in '-'")
	}

	if !s.Accept(".") {
		return fmt.Errorf("did not find a '.'")
	}

	s.scanningHostname = true

	if !IsAlphaNum(s.Peek()) {
		// got invalid character for hostname (either - in first char, or other characters)
		s.Emit(ItemHostname)
	}

	return nil
}

// ScanHostname scans & emits a valid hostname
func ScanHostname(s *Scanr) StateFn {
	if err := getHostnamePart(s); err != nil {
		s.Emit(ItemError)
	}
	s.Emit(ItemHostname)

	return s.homeState
}
