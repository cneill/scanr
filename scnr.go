package scnr

import (
	"strings"
	"unicode/utf8"
)

// StateFn is a function used to transition between states in the Scanner
type StateFn func(*Scanner) StateFn

// RuneFn is a function used to determine if a rune is in a given set of characters
type RuneFn func(rune) bool

// ItemType tells the parser what kind of Item this is
type ItemType int

// Item represents a string of a particular ItemType
type Item struct {
	typ ItemType
	pos int
	val string
}

// Items is a convenience type for a list of Item
type Items []Item

// Scanner is used to emit items / tokens to be parsed by higher-level logic
type Scanner struct {
	input    string
	state    StateFn
	pos      int
	start    int
	width    int
	lastPos  int
	items    chan Item
	lastItem Item
}

// NewScanner returns an initialized Scanner
func NewScanner(start StateFn) *Scanner {
	s := &Scanner{
		state: start,
		items: make(chan Item),
	}
	return s
}

// Run kicks off the scanner, starting with the first StateFn
func (s *Scanner) Run(input string) {
	s.input = input
	for s.state != nil {
		s.state = s.state(s)
	}
}

// next returns the next rune in the input.
func (s *Scanner) next() rune {
	if s.pos >= len(s.input) {
		s.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(s.input[s.pos:])
	s.width = w
	s.pos += s.width
	return r
}

// backup steps back one rone. Can only be called once per call of next.
func (s *Scanner) backup() {
	s.pos -= s.width
}

// peek returns but does not consume the next rune in the input.
func (s *Scanner) peek() rune {
	r := s.next()
	s.backup()
	return r
}

// emit passes an item to the items channel.
func (s *Scanner) emit(t ItemType) {
	i := Item{t, s.start, s.input[s.start:s.pos]}
	s.items <- i
	s.lastItem = i
	s.start = s.pos
}

// accept consumes the next rune if it's from the valid set
func (s *Scanner) accept(valid string) bool {
	if strings.ContainsRune(valid, s.next()) {
		return true
	}
	s.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set
func (s *Scanner) acceptRun(valid string) int {
	var length = 0
	for strings.ContainsRune(valid, s.next()) {
		length++
	}
	s.backup()
	return length
}

// acceptUntil consumes to 'end' or eof; returns true if it accepts, false otherwise
func (s *Scanner) acceptUntil(end rune) bool {
	if s.peek() == end || s.peek() == eof {
		return false
	}
	for r := s.next(); r != end && r != eof; r = s.next() {
	}
	s.backup()
	return true
}

// acceptWhileRuneFn consumes while 'fn' returns true
func (s *Scanner) acceptWhileRuneFn(fn RuneFn) bool {
	accepted := false
	for r := s.next(); fn(r) && r != eof; r = s.next() {
		accepted = true
	}
	s.backup()
	return accepted
}

// acceptUntilRuneFn consumes until 'end' returns true
func (s *Scanner) acceptUntilRuneFn(end RuneFn) bool {
	accepted := false
	for r := s.next(); !end(r) && r != eof; r = s.next() {
		// for r := s.peek(); !end(r) && r != eof; r = s.peek() {
		accepted = true
	}
	s.backup()
	return accepted
}

// acceptSequence consumes a string if found & returns true, false if not
func (s *Scanner) acceptSequence(valid string) bool {
	if strings.HasPrefix(s.input[s.pos:], valid) {
		s.pos += len(valid)
		return true
	}
	return false
}

// nextItem returns the next Item from the input; called by parser
func (s *Scanner) nextItem() Item {
	item := <-s.items
	s.lastPos = item.pos
	return item
}

// drain runs through output so lexing goroutine exists; called by parser
func (s *Scanner) drain() {
	for range s.items {
	}
}

// ignore skips over the pending input before this point. - UNUSED FOR NOW
func (s *Scanner) ignore() {
	s.start = s.pos
}

// IsQuote returns true if r is one of " ' `
func (s *Scanner) IsQuote(r rune) bool {
	return strings.ContainsRune("\"'`", r)
}

// IsNewline returns ture if r is one of \r or \n
func (s *Scanner) IsNewline(r rune) bool {
	return r == '\r' || r == '\n'
}

// IsAlphaLower returns true if r is between a and z
func (s *Scanner) IsAlphaLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

// IsAlphaUpper returns true if r is between A and Z
func (s *Scanner) IsAlphaUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// IsAlpha returns true if r is between a and z or A and Z
func (s *Scanner) IsAlpha(r rune) bool {
	return s.IsAlphaLower(r) || s.IsAlphaUpper(r)
}

// IsNumber returns true if r is between 0 and 9
func (s *Scanner) IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsAlphaNum returns true if r is a letter or number
func (s *Scanner) IsAlphaNum(r rune) bool {
	return s.IsAlphaLower(r) || s.IsAlphaUpper(r) || s.IsNumber(r)
}
