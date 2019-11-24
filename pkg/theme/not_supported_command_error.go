package theme

import (
	"fmt"
	"strings"
)

type NotSupportedCommandError struct {
	Command           string
	SupportedCommands []string
}

func (s *NotSupportedCommandError) Error() string {
	return fmt.Sprintf(
		"%s is not one of the supported commands: %s",
		s.Command,
		strings.Join(s.SupportedCommands, ", "),
	)
}
