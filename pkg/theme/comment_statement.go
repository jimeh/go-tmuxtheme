package theme

import (
	"regexp"
	"strings"
)

var commandStatementCommands = []string{"#"}
var commentStatementMatcher = regexp.MustCompile(`^\s*(#)(.*?)$`)

type CommentStatement struct {
	Msg string
}

func (s *CommentStatement) Parse(body string) error {
	match := commentStatementMatcher.FindStringSubmatch(body)
	if len(match) < 3 {
		return &NotSupportedCommandError{
			strings.SplitN(strings.TrimSpace(body), " ", 2)[0],
			commandStatementCommands,
		}
	}

	s.Msg = strings.TrimSpace(match[2])
	return nil
}

func (s *CommentStatement) Execute(theme *Theme) error {
	return nil
}
