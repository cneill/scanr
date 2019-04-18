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

// Next returns the next rune in the input.
func (s *Scanner) Next() rune {
	if s.pos >= len(s.input) {
		s.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(s.input[s.pos:])
	s.width = w
	s.pos += s.width
	return r
}

// Backup steps back one rune. Can only be called once per call of next.
func (s *Scanner) Backup() {
	s.pos -= s.width
}

// Peek returns but does not consume the next rune in the input.
func (s *Scanner) Peek() rune {
	r := s.Next()
	s.Backup()
	return r
}

// Emit passes an item to the items channel.
func (s *Scanner) Emit(t ItemType) {
	i := Item{t, s.start, s.input[s.start:s.pos]}
	s.items <- i
	s.lastItem = i
	s.start = s.pos
}

// Accept consumes the next rune if it's from the valid set
func (s *Scanner) Accept(valid string) bool {
	if strings.ContainsRune(valid, s.Next()) {
		return true
	}
	s.Backup()
	return false
}

// AcceptRun consumes a run of runes from the valid set
func (s *Scanner) AcceptRun(valid string) int {
	var length = 0
	for strings.ContainsRune(valid, s.Next()) {
		length++
	}
	s.Backup()
	return length
}

// AcceptUntil consumes to 'end' or eof; returns true if it accepts, false otherwise
func (s *Scanner) AcceptUntil(end rune) bool {
	if s.Peek() == end || s.Peek() == eof {
		return false
	}
	for r := s.Next(); r != end && r != eof; r = s.Next() {
	}
	s.Backup()
	return true
}

// AcceptWhileRuneFn consumes while 'fn' returns true
func (s *Scanner) AcceptWhileRuneFn(fn RuneFn) bool {
	accepted := false
	for r := s.Next(); fn(r) && r != eof; r = s.Next() {
		accepted = true
	}
	s.Backup()
	return accepted
}

// AcceptUntilRuneFn consumes until 'end' returns true
func (s *Scanner) AcceptUntilRuneFn(end RuneFn) bool {
	accepted := false
	for r := s.Next(); !end(r) && r != eof; r = s.Next() {
		// for r := s.peek(); !end(r) && r != eof; r = s.peek() {
		accepted = true
	}
	s.Backup()
	return accepted
}

// AcceptSequence consumes a string if found & returns true, false if not
func (s *Scanner) AcceptSequence(valid string) bool {
	if strings.HasPrefix(s.input[s.pos:], valid) {
		s.pos += len(valid)
		return true
	}
	return false
}

// NextItem returns the next Item from the input; called by parser
func (s *Scanner) NextItem() Item {
	item := <-s.items
	s.lastPos = item.pos
	return item
}

// Drain runs through output so lexing goroutine exists; called by parser
func (s *Scanner) Drain() {
	for range s.items {
	}
}

// Ignore skips over the pending input before this point. - UNUSED FOR NOW
func (s *Scanner) Ignore() {
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
