package theme

type ArgumentError struct {
	msg string
}

func (s *ArgumentError) Error() string {
	return s.msg
}
