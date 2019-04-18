package scanr

// ScanSpace scans & emits a run of space characters
func ScanSpace(s *Scanner) StateFn {
	s.AcceptWhileRuneFn(s.IsSpace)
	s.Emit(ItemSpace)
	return s.homeState
}

// ScanNewline scans & emits a Newline of the variants "\n", "\r", and "\r\n"
func ScanNewline(s *Scanner) StateFn {
	n := s.Next()
	if !s.IsNewline(n) {
		s.Backup()
		return s.homeState
	} else if n == '\r' {
		s.Accept("\n")
	}

	s.Emit(ItemNewline)
	return s.homeState
}
