package scanr

const (
	ItemError   ItemType = iota // 0
	ItemEOF                     // 1
	ItemSpace                   // 2
	ItemNewline                 // 3
)

const (
	eof = rune(0)
)
