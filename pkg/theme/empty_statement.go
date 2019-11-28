package theme

import "github.com/kballard/go-shellquote"

type EmptyStatement struct {
}

func (s *EmptyStatement) Parse(body string) error {
	args, err := shellquote.Split(body)
	if err != nil {
		return nil
	} else if len(args) > 0 {
		return &NotSupportedCommandError{"", []string{}}
	}

	return nil
}

func (s *EmptyStatement) Execute(theme *Theme) error {
	return nil
}
