package theme

import (
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/kballard/go-shellquote"
)

var setOptionStatementCommands = []string{
	"set", "set-option", "set-window-option",
}
var setOptionStatementFormatPattern = regexp.MustCompile(
	"#{(@?[a-zA-Z0-9-_]+?)}",
)

type SetOptionFlags struct {
	Append      bool   `short:"a"`
	Format      bool   `short:"F"`
	Global      bool   `short:"g"`
	OnlyIfUnset bool   `short:"o"`
	Quiet       bool   `short:"q"`
	Server      bool   `short:"s"`
	Target      string `short:"t"`
	Unset       bool   `short:"u"`
	Window      bool   `short:"w"`
}

type SetOptionStatement struct {
	Flags  *SetOptionFlags
	Option string
	Value  string
}

func (s *SetOptionStatement) Parse(body string) error {
	args, err := shellquote.Split(body)
	if err != nil {
		return &NotSupportedCommandError{
			strings.SplitN(strings.TrimSpace(body), " ", 2)[0],
			setOptionStatementCommands,
		}
	}

	args, err = s.parseCommand(args)
	if err != nil {
		return err
	}

	args, err = s.parseFlags(args)
	if err != nil {
		return err
	}

	return s.parseArguments(args)
}

func (s *SetOptionStatement) Execute(theme *Theme) error {
	if s.Flags.Server {
		return s.applyValue(theme, theme.ServerOptions)
	} else if s.Flags.Global && s.Flags.Window {
		return s.applyValue(theme, theme.GlobalWindowOptions)
	} else if s.Flags.Window {
		return s.applyValue(theme, theme.WindowOptions)
	} else if s.Flags.Global {
		return s.applyValue(theme, theme.GlobalSessionOptions)
	} else {
		return s.applyValue(theme, theme.SessionOptions)
	}
}

func (s *SetOptionStatement) parseCommand(args []string) ([]string, error) {
	cmd := ""

	if len(args) > 1 {
		cmd, args = args[0], args[1:]
		for _, c := range setOptionStatementCommands {
			if cmd == c {
				if cmd == "set-window-option" {
					args = append([]string{"-w"}, args...)
				}

				return args, nil
			}
		}
	} else {
		if len(args) == 1 {
			cmd = args[0]
		}
		args = []string{}
	}

	return args, &NotSupportedCommandError{cmd, setOptionStatementCommands}
}

func (s *SetOptionStatement) parseFlags(args []string) ([]string, error) {
	s.Flags = &SetOptionFlags{}
	args, err := flags.ParseArgs(s.Flags, args)
	if err != nil {
		return nil, err
	}

	return args, nil
}

func (s *SetOptionStatement) parseArguments(args []string) error {
	if len(args) == 0 {
		return &NoOptionArgumentError{}
	}

	s.Option = args[0]
	if len(args) > 1 {
		s.Value = args[1]
	}
	return nil
}

func (s *SetOptionStatement) applyValue(theme *Theme, options map[string]string) error {
	option := s.Option
	value := s.Value

	if s.Flags.OnlyIfUnset {
		if _, ok := options[option]; ok {
			return nil
		}
	}

	if s.Flags.Unset {
		delete(options, option)
		return nil
	}

	if s.Flags.Format {
		value = s.formatValue(theme, value)
	}

	if s.Flags.Append {
		options[option] = options[option] + value
	} else {
		options[option] = value
	}

	return nil
}

func (s *SetOptionStatement) formatValue(theme *Theme, value string) string {
	return setOptionStatementFormatPattern.ReplaceAllStringFunc(
		value,
		func(match string) string {
			name := setOptionStatementFormatPattern.ReplaceAllString(
				match, `$1`,
			)
			return s.lookupOptionValue(theme, name)
		},
	)
}

func (s *SetOptionStatement) lookupOptionValue(theme *Theme, name string) string {
	if val, ok := theme.WindowOptions[name]; ok {
		return val
	} else if val, ok := theme.GlobalWindowOptions[name]; ok {
		return val
	} else if val, ok := theme.SessionOptions[name]; ok {
		return val
	} else if val, ok := theme.GlobalSessionOptions[name]; ok {
		return val
	} else if val, ok := theme.ServerOptions[name]; ok {
		return val
	}

	return ""
}
