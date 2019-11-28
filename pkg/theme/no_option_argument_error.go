package theme

type NoOptionArgumentError struct{}

func (s *NoOptionArgumentError) Error() string {
	return "No option argument given"
}
