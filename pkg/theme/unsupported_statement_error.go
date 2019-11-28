package theme

import "fmt"

type UnsupportedStatementError struct {
	Body string
}

func (s *UnsupportedStatementError) Error() string {
	return fmt.Sprintf("Unsupported statement: %s", s.Body)
}
