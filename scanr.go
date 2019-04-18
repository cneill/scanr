package scanr

import (
	"strings"
	"unicode/utf8"
)

// StateFn is a function used to transition between states in the Scanr
type StateFn func(*Scanr) StateFn

// RuneFn is a function used to determine if a rune is in a given set of characters
type RuneFn func(rune) bool

// ItemType tells the parser what kind of Item this is
type ItemType int

// Item represents a string of a particular ItemType
type Item struct {
	Typ ItemType
	Pos int
	Val string
}

// Items is a convenience type for a list of Item
type Items []Item

// Scanr is used to emit items / tokens to be parsed by higher-level logic
type Scanr struct {
	input     string
	homeState StateFn
	state     StateFn
	pos       int
	start     int
	width     int
	lastPos   int
	items     chan Item
	lastItem  Item
}

// NewScanr returns an initialized Scanr
func NewScanr(home StateFn) *Scanr {
	s := &Scanr{
		homeState: home,
		state:     home,
		items:     make(chan Item),
	}
	return s
}

// Run kicks off the scanner, starting with the first StateFn
func (s *Scanr) Run(input string) {
	s.input = input
	for s.state != nil {
		s.state = s.state(s)
	}
}

// Next returns the next rune in the input.
func (s *Scanr) Next() rune {
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
func (s *Scanr) Backup() {
	s.pos -= s.width
}

// Peek returns but does not consume the next rune in the input.
func (s *Scanr) Peek() rune {
	r := s.Next()
	s.Backup()
	return r
}

// Emit passes an item to the items channel.
func (s *Scanr) Emit(t ItemType) {
	i := Item{t, s.start, s.input[s.start:s.pos]}
	s.items <- i
	s.lastItem = i
	s.start = s.pos
}

// Accept consumes the next rune if it's from the valid set
func (s *Scanr) Accept(valid string) bool {
	if strings.ContainsRune(valid, s.Next()) {
		return true
	}
	s.Backup()
	return false
}

// AcceptRun consumes a run of runes from the valid set
func (s *Scanr) AcceptRun(valid string) int {
	var length = 0
	for strings.ContainsRune(valid, s.Next()) {
		length++
	}
	s.Backup()
	return length
}

// AcceptUntil consumes to 'end' or eof; returns true if it accepts, false otherwise
func (s *Scanr) AcceptUntil(end rune) bool {
	if s.Peek() == end || s.Peek() == eof {
		return false
	}
	for r := s.Next(); r != end && r != eof; r = s.Next() {
	}
	s.Backup()
	return true
}

// AcceptWhileRuneFn consumes while 'fn' returns true
func (s *Scanr) AcceptWhileRuneFn(fn RuneFn) bool {
	accepted := false
	for r := s.Next(); fn(r) && r != eof; r = s.Next() {
		accepted = true
	}
	s.Backup()
	return accepted
}

// AcceptUntilRuneFn consumes until 'end' returns true
func (s *Scanr) AcceptUntilRuneFn(end RuneFn) bool {
	accepted := false
	for r := s.Next(); !end(r) && r != eof; r = s.Next() {
		// for r := s.peek(); !end(r) && r != eof; r = s.peek() {
		accepted = true
	}
	s.Backup()
	return accepted
}

// AcceptSequence consumes a string if found & returns true, false if not
func (s *Scanr) AcceptSequence(valid string) bool {
	if strings.HasPrefix(s.input[s.pos:], valid) {
		s.pos += len(valid)
		return true
	}
	return false
}

// NextItem returns the next Item from the input; called by parser
func (s *Scanr) NextItem() Item {
	item := <-s.items
	s.lastPos = item.Pos
	return item
}

// Drain runs through output so lexing goroutine exists; called by parser
func (s *Scanr) Drain() {
	for range s.items {
	}
}

// Ignore skips over the pending input before this point. - UNUSED FOR NOW
func (s *Scanr) Ignore() {
	s.start = s.pos
}

// IsSpace returns true if r is a space character
func (s *Scanr) IsSpace(r rune) bool {
	return r == ' '
}

// IsWhitespace returns true if r is a space character or tab character
func (s *Scanr) IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

// IsQuote returns true if r is one of " ' `
func (s *Scanr) IsQuote(r rune) bool {
	return strings.ContainsRune("\"'`", r)
}

// IsNewline returns ture if r is one of \r or \n
func (s *Scanr) IsNewline(r rune) bool {
	return r == '\r' || r == '\n'
}

// IsAlphaLower returns true if r is between a and z
func (s *Scanr) IsAlphaLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

// IsAlphaUpper returns true if r is between A and Z
func (s *Scanr) IsAlphaUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// IsAlpha returns true if r is between a and z or A and Z
func (s *Scanr) IsAlpha(r rune) bool {
	return s.IsAlphaLower(r) || s.IsAlphaUpper(r)
}

// IsNumber returns true if r is between 0 and 9
func (s *Scanr) IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsAlphaNum returns true if r is a letter or number
func (s *Scanr) IsAlphaNum(r rune) bool {
	return s.IsAlphaLower(r) || s.IsAlphaUpper(r) || s.IsNumber(r)
}
