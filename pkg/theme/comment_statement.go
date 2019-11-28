package theme

import (
	"regexp"
	"strings"

	"github.com/kballard/go-shellquote"
)

var commandStatementCommands = []string{"#"}
var commentStatementMatcher = regexp.MustCompile(`^\s*(#)(.*?)$`)

type CommentStatement struct {
	Msg string
}

func (s *CommentStatement) Parse(body string) error {
	args, err := shellquote.Split(body)
	if err != nil {
		return err
	} else if len(args) == 0 {
		return &NotSupportedCommandError{"", commandStatementCommands}
	}

	match := commentStatementMatcher.FindStringSubmatch(args[0])
	if len(match) < 3 {
		return &NotSupportedCommandError{args[0], commandStatementCommands}
	}

	err = s.parseCommand(match[1])
	if err != nil {
		return err
	}

	msg := []string{}
	if match[2] != "" {
		msg = append(msg, match[2])
	}
	msg = append(msg, args[1:]...)

	s.Msg = strings.Join(msg, " ")
	return nil
}

func (s *CommentStatement) Execute(theme *Theme) error {
	return nil
}

func (s *CommentStatement) parseCommand(cmd string) error {
	for _, supportedCommand := range commandStatementCommands {
		if cmd == supportedCommand {
			return nil
		}
	}

	return &NotSupportedCommandError{cmd, commandStatementCommands}
}
