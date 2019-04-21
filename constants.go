package scanr

const (
	ItemError    ItemType = iota // 0
	ItemEOF                      // 1
	ItemSpace                    // 2
	ItemNewline                  // 3
	ItemIP                       // 4
	ItemHostname                 // 5
)

const (
	eof    = rune(0)
	digits = "1234567890"
)
