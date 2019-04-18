package scanr

// ScanSpace scans & emits a run of space characters
func ScanSpace(s *Scanner) StateFn {
	s.AcceptWhileRuneFn(s.IsSpace)
	s.Emit(ItemSpace)
	return s.homeState
}

// ScanNewline scans & emits a Newline of the variants "\n", "\r", and "\r\n"
func ScanNewline(s *Scanner) StateFn {
	if n := s.Next(); n == '\r' {
		s.Accept("\n")
	} else if n != '\n' {
		s.Backup()
	}
	s.Emit(ItemNewline)
	return s.homeState
}
