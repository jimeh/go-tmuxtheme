package theme

import (
	"regexp"

	"github.com/jessevdk/go-flags"
	"github.com/kballard/go-shellquote"
)

var setOptionCommands = []string{"set", "set-option", "set-window-option"}
var setOptionFormatPattern = regexp.MustCompile("#{(@?[a-zA-Z0-9-_]+?)}")

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
	Flags     *SetOptionFlags
	Arguments []string
	theme     *Theme
}

func (s *SetOptionStatement) Parse(body string) error {
	parts, err := shellquote.Split(body)
	if err != nil {
		return err
	}
	cmd, args := parts[0], parts[1:]

	err = s.validateCommand(cmd)
	if err != nil {
		return err
	}

	s.Flags = &SetOptionFlags{}
	s.Arguments = []string{}

	if cmd == "set-window-option" {
		args = append([]string{"-w"}, args...)
	}

	s.Arguments, err = flags.ParseArgs(s.Flags, args)
	if err != nil {
		return err
	}

	return nil
}

func (s *SetOptionStatement) Execute() error {
	if s.Flags.Server {
		return s.applyValue(&s.theme.ServerOptions)
	} else if s.Flags.Global && s.Flags.Window {
		return s.applyValue(&s.theme.GlobalWindowOptions)
	} else if s.Flags.Window {
		return s.applyValue(&s.theme.WindowOptions)
	} else if s.Flags.Global {
		return s.applyValue(&s.theme.GlobalSessionOptions)
	} else {
		return s.applyValue(&s.theme.SessionOptions)
	}
}

func (s *SetOptionStatement) validateCommand(cmd string) error {
	for _, supportedCommand := range setOptionCommands {
		if cmd == supportedCommand {
			return nil
		}
	}

	return &NotSupportedCommandError{cmd, setOptionCommands}
}

func (s *SetOptionStatement) applyValue(options *map[string]string) error {
	name := s.Arguments[0]

	if s.Flags.OnlyIfUnset {
		if _, ok := (*options)[name]; ok {
			return nil
		}
	}

	if s.Flags.Unset {
		delete((*options), name)
		return nil
	}

	value := s.Arguments[1]
	if s.Flags.Format {
		value = s.formatValue(value)
	}

	if s.Flags.Append {
		(*options)[name] = (*options)[name] + value
	} else {
		(*options)[name] = value
	}

	return nil
}

func (s *SetOptionStatement) formatValue(value string) string {
	return setOptionFormatPattern.ReplaceAllStringFunc(
		value,
		func(match string) string {
			name := setOptionFormatPattern.ReplaceAllString(match, `$1`)
			return s.lookupOptionValue(name)
		},
	)
}

func (s *SetOptionStatement) lookupOptionValue(name string) string {
	if val, ok := s.theme.WindowOptions[name]; ok {
		return val
	} else if val, ok := s.theme.GlobalWindowOptions[name]; ok {
		return val
	} else if val, ok := s.theme.SessionOptions[name]; ok {
		return val
	} else if val, ok := s.theme.GlobalSessionOptions[name]; ok {
		return val
	} else if val, ok := s.theme.ServerOptions[name]; ok {
		return val
	}

	return ""
}
