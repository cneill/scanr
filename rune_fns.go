package scanr

import "strings"

// RuneFn is a function used to determine if a rune is in a given set of characters
type RuneFn func(rune) bool

// IsSpace returns true if r is a space character
func IsSpace(r rune) bool {
	return r == ' '
}

// IsWhitespace returns true if r is a space character or tab character
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

// IsQuote returns true if r is one of " ' `
func IsQuote(r rune) bool {
	return strings.ContainsRune("\"'`", r)
}

// IsNewline returns ture if r is one of \r or \n
func IsNewline(r rune) bool {
	return r == '\r' || r == '\n'
}

// IsAlphaLower returns true if r is between a and z
func IsAlphaLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

// IsAlphaUpper returns true if r is between A and Z
func IsAlphaUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// IsAlpha returns true if r is between a and z or A and Z
func IsAlpha(r rune) bool {
	return IsAlphaLower(r) || IsAlphaUpper(r)
}

// IsNumber returns true if r is between 0 and 9
func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsAlphaNum returns true if r is a letter or number
func IsAlphaNum(r rune) bool {
	return IsAlphaLower(r) || IsAlphaUpper(r) || IsNumber(r)
}

// IsHostnameChar returns true if r is a letter, number, or "-"
func IsHostnameChar(r rune) bool {
	return IsAlphaNum(r) || r == '-'
}
