package scanr

import (
	"strings"
	"unicode/utf8"
)

/*
TODO:
- allow different state & emission paths for state functions - a generic "parse hostname" state fn might be used to JUST get that hostname and emit it; it might also be part of a URL parsing flow
	- using state fns makes this hard - can't pass params easily?
	- maybe have wrapper state fns, where specifying a "next" and "emit/noemit" is possible?
	- without doing something like this ^ all state fns must be unique, even if they do the same things. multiple hostname parser functions, etc.
	- if structs are used, it makes each "state fn" more cumbersome to construct, but gives more configurability; can include a function with bools for emit/noemit, perhaps a whole state chain with
	  a custom emit type?
    - I had "chains" of states like this in axe - perhaps look at that approach again? I didn't allow emit configurability there, merely ordering, so will have to figure that out
- figure out how to test state fns - need to be able to do the Parse(), NextItem(), etc. bits all within a test
- add more general-purpose state function building blocks that can be combined and rearranged for different purposes:
	- "number" - i.e. int/float/imaginary/etc
	- "word" - i.e. set of letters, separated by spaces
	- URL
	- quoted string
*/

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

	// TODO: wtf
	scanningHostname bool
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

// PrevRune returns the last rune that was scanned
func (s *Scanr) PrevRune() rune {
	r, _ := utf8.DecodeRuneInString(s.input[s.pos-s.width:])
	return r
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
