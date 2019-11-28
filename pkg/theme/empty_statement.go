package theme

import "strings"

type EmptyStatement struct {
}

func (s *EmptyStatement) Parse(body string) error {
	trimmed := strings.TrimSpace(body)

	if trimmed != "" {
		return &NotSupportedCommandError{
			strings.SplitN(trimmed, " ", 2)[0],
			[]string{},
		}
	}

	return nil
}

func (s *EmptyStatement) Execute(theme *Theme) error {
	return nil
}
