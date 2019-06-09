package scanr

// ItemType tells the parser what kind of Item this is
type ItemType int

// Item represents a string of a particular ItemType
type Item struct {
	Typ ItemType
	Pos int
	Val string
}

// Items is a convenience type for an array of Item
type Items []Item

// String returns the concatenated values of all Item in this array
func (i Items) String() string {
	var result = ""
	for _, item := range i {
		result += item.Val
	}
	return result
}
