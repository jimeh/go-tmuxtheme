package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetOptionStatementInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*Statement)(nil), &SetOptionStatement{})
}

func TestSetOptionStatementParse(t *testing.T) {
	var tests = []struct {
		body   string
		flags  *SetOptionFlags
		args   []string
		option string
		value  string
		error  error
	}{
		{
			body:   `set -a myopt foo`,
			flags:  &SetOptionFlags{Append: true},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set -F myopt 'foo bar'`,
			flags:  &SetOptionFlags{Format: true},
			option: "myopt",
			value:  "foo bar",
		},
		{
			body:   `set -F myopt ' foo bar  '`,
			flags:  &SetOptionFlags{Format: true},
			option: "myopt",
			value:  " foo bar  ",
		},
		{
			body:   `set -g myopt "foo bar"`,
			flags:  &SetOptionFlags{Global: true},
			option: "myopt",
			value:  "foo bar",
		},
		{
			body:   `set -g myopt "  foo bar "`,
			flags:  &SetOptionFlags{Global: true},
			option: "myopt",
			value:  "  foo bar  ",
		},
		{
			body:   `set -o myopt foo`,
			flags:  &SetOptionFlags{OnlyIfUnset: true},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set -q myopt foo`,
			flags:  &SetOptionFlags{Quiet: true},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set -s myopt foo`,
			flags:  &SetOptionFlags{Server: true},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set -t other:3 myopt foo`,
			flags:  &SetOptionFlags{Target: "other:3"},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set -u myopt`,
			flags:  &SetOptionFlags{Unset: true},
			option: "myopt",
		},
		{
			body:   `set -w myopt foo`,
			flags:  &SetOptionFlags{Window: true},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set-option -w myopt foo`,
			flags:  &SetOptionFlags{Window: true},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set-window-option myopt foo`,
			flags:  &SetOptionFlags{Window: true},
			option: "myopt",
			value:  "foo",
		},
		{
			body:   `set -goq @myopt 'hello world'`,
			flags:  &SetOptionFlags{Global: true, OnlyIfUnset: true, Quiet: true},
			option: "@myopt",
			value:  "hello world",
		},
		{
			body:   `set -gF @myopt 'hello #{@other} world'`,
			flags:  &SetOptionFlags{Global: true, Format: true},
			option: "@myopt",
			value:  "hello #{@other} world",
		},
		{
			body: `has-session -t myopt`,
			error: &NotSupportedCommandError{
				"has-session", setOptionStatementCommands,
			},
		},
		{
			body:  ``,
			error: &NotSupportedCommandError{"", setOptionStatementCommands},
		},
		{
			body:  `set`,
			error: &NotSupportedCommandError{"set", setOptionStatementCommands},
		},
		{
			body:  `set -gu`,
			error: &NoOptionArgumentError{},
		},
	}

	for _, tt := range tests {
		s := &SetOptionStatement{}

		err := s.Parse(tt.body)

		if tt.flags != nil {
			assert.Equal(t, tt.flags, s.Flags)
		}

		if tt.option != "" {
			assert.Equal(t, tt.option, s.Option)
		}

		if tt.error != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.error, err)
		}
	}
}

func TestSetOptionStatementExecute(t *testing.T) {
	var tests = []struct {
		body               string
		server             map[string]string
		globalSession      map[string]string
		session            map[string]string
		globalWindow       map[string]string
		window             map[string]string
		serverSetup        map[string]string
		globalSessionSetup map[string]string
		sessionSetup       map[string]string
		globalWindowSetup  map[string]string
		windowSetup        map[string]string
	}{
		//
		// Session Options
		//
		{
			body:    `set @name "John"`,
			session: map[string]string{"@name": "John"},
		},
		{
			body:    `set-option @name "John"`,
			session: map[string]string{"@name": "John"},
		},
		{
			body:    `set -o @name "Jim"`,
			session: map[string]string{"@name": "Jim"},
		},
		{
			body:         `set -o @name "Jim"`,
			sessionSetup: map[string]string{"@name": "John"},
			session:      map[string]string{"@name": "John"},
		},
		{
			body:         `set -a @name "Jim"`,
			sessionSetup: map[string]string{"@name": "John"},
			session:      map[string]string{"@name": "JohnJim"},
		},
		{
			body:         `set -u @name`,
			sessionSetup: map[string]string{"@name": "John", "@foo": "bar"},
			session:      map[string]string{"@foo": "bar"},
		},
		{
			body:         `set -F @message "Hi #{@name}"`,
			sessionSetup: map[string]string{"@name": "John Smith"},
			session: map[string]string{
				"@name":    "John Smith",
				"@message": "Hi John Smith",
			},
		},
		//
		// Server Options
		//
		{
			body:   `set -s @name "John"`,
			server: map[string]string{"@name": "John"},
		},
		{
			body:   `set-option -s @name "John"`,
			server: map[string]string{"@name": "John"},
		},
		{
			body:   `set -so @name "Jim"`,
			server: map[string]string{"@name": "Jim"},
		},
		{
			body:        `set -so @name "Jim"`,
			serverSetup: map[string]string{"@name": "John"},
			server:      map[string]string{"@name": "John"},
		},
		{
			body:        `set -sa @name "Jim"`,
			serverSetup: map[string]string{"@name": "John"},
			server:      map[string]string{"@name": "JohnJim"},
		},
		{
			body:        `set -su @name`,
			serverSetup: map[string]string{"@name": "John", "@foo": "bar"},
			server:      map[string]string{"@foo": "bar"},
		},
		{
			body:        `set -sF @message "Hi #{@name}"`,
			serverSetup: map[string]string{"@name": "John Smith"},
			server: map[string]string{
				"@name":    "John Smith",
				"@message": "Hi John Smith",
			},
		},
		//
		// Global Session Options
		//
		{
			body:          `set -g @name "John"`,
			globalSession: map[string]string{"@name": "John"},
		},
		{
			body:          `set-option -g @name "John"`,
			globalSession: map[string]string{"@name": "John"},
		},
		{
			body:          `set -go @name "Jim"`,
			globalSession: map[string]string{"@name": "Jim"},
		},
		{
			body:               `set -go @name "Jim"`,
			globalSessionSetup: map[string]string{"@name": "John"},
			globalSession:      map[string]string{"@name": "John"},
		},
		{
			body:               `set -ga @name "Jim"`,
			globalSessionSetup: map[string]string{"@name": "John"},
			globalSession:      map[string]string{"@name": "JohnJim"},
		},
		{
			body:               `set -gu @name`,
			globalSessionSetup: map[string]string{"@name": "John", "@foo": "bar"},
			globalSession:      map[string]string{"@foo": "bar"},
		},
		{
			body:               `set -gF @message "Hi #{@name}"`,
			globalSessionSetup: map[string]string{"@name": "John Smith"},
			globalSession: map[string]string{
				"@name":    "John Smith",
				"@message": "Hi John Smith",
			},
		},
		//
		// Window Options
		//
		{
			body:   `set -w @name "John"`,
			window: map[string]string{"@name": "John"},
		},
		{
			body:   `set-option -w @name "John"`,
			window: map[string]string{"@name": "John"},
		},
		{
			body:   `set-window-option @name "John"`,
			window: map[string]string{"@name": "John"},
		},
		{
			body:   `set -wo @name "Jim"`,
			window: map[string]string{"@name": "Jim"},
		},
		{
			body:        `set -wo @name "Jim"`,
			windowSetup: map[string]string{"@name": "John"},
			window:      map[string]string{"@name": "John"},
		},
		{
			body:        `set -wa @name "Jim"`,
			windowSetup: map[string]string{"@name": "John"},
			window:      map[string]string{"@name": "JohnJim"},
		},
		{
			body:        `set -wu @name`,
			windowSetup: map[string]string{"@name": "John", "@foo": "bar"},
			window:      map[string]string{"@foo": "bar"},
		},
		{
			body:        `set -wF @message "Hi #{@name}"`,
			windowSetup: map[string]string{"@name": "John Smith"},
			window: map[string]string{
				"@name":    "John Smith",
				"@message": "Hi John Smith",
			},
		},
		//
		// Global Window Options
		//
		{
			body:         `set -gw @name "John"`,
			globalWindow: map[string]string{"@name": "John"},
		},
		{
			body:         `set-option -wg @name "John"`,
			globalWindow: map[string]string{"@name": "John"},
		},
		{
			body:         `set-window-option -g @name "John"`,
			globalWindow: map[string]string{"@name": "John"},
		},
		{
			body:         `set -wgo @name "Jim"`,
			globalWindow: map[string]string{"@name": "Jim"},
		},
		{
			body:              `set -gwo @name "Jim"`,
			globalWindowSetup: map[string]string{"@name": "John"},
			globalWindow:      map[string]string{"@name": "John"},
		},
		{
			body:              `set -wga @name "Jim"`,
			globalWindowSetup: map[string]string{"@name": "John"},
			globalWindow:      map[string]string{"@name": "JohnJim"},
		},
		{
			body:              `set -gwu @name`,
			globalWindowSetup: map[string]string{"@name": "John", "@foo": "bar"},
			globalWindow:      map[string]string{"@foo": "bar"},
		},
		{
			body:              `set -wgF @message "Hi #{@name}"`,
			globalWindowSetup: map[string]string{"@name": "John Smith"},
			globalWindow: map[string]string{
				"@name":    "John Smith",
				"@message": "Hi John Smith",
			},
		},
		//
		// Formatting
		//
		{
			body:         `set -F @foo "foo #{@bar} baz"`,
			sessionSetup: map[string]string{"@bar": "bar"},
			session:      map[string]string{"@bar": "bar", "@foo": "foo bar baz"},
		},
		{
			body:         `set -F @foo "foo #{bar} baz"`,
			sessionSetup: map[string]string{"bar": "bar"},
			session:      map[string]string{"bar": "bar", "@foo": "foo bar baz"},
		},
		{
			body:         `set -F @foo "foo #{@bar}#{@bar}"`,
			sessionSetup: map[string]string{"@bar": "bar"},
			session:      map[string]string{"@bar": "bar", "@foo": "foo barbar"},
		},
		{
			body:         `set -F @msg "#{@foo} #{@bar}"`,
			sessionSetup: map[string]string{"@foo": "foo", "@bar": "bar"},
			session: map[string]string{
				"@foo": "foo",
				"@bar": "bar",
				"@msg": "foo bar",
			},
		},
		{
			body:               `set -gF @foo "foo #{@bar} baz"`,
			globalSessionSetup: map[string]string{"@bar": "bar"},
			globalSession: map[string]string{
				"@bar": "bar",
				"@foo": "foo bar baz",
			},
		},
		{
			body:         `set -Fo @foo "foo #{@bar} baz"`,
			sessionSetup: map[string]string{"@bar": "bar", "@foo": "foo"},
			session:      map[string]string{"@bar": "bar", "@foo": "foo"},
		},
		{
			body:         `set -Fa @foo " #{@bar} baz"`,
			sessionSetup: map[string]string{"@bar": "bar", "@foo": "foo"},
			session:      map[string]string{"@bar": "bar", "@foo": "foo bar baz"},
		},
		{
			body:          `set -gF @message "Hi #{@name}"`,
			windowSetup:   map[string]string{"@name": "John"},
			window:        map[string]string{"@name": "John"},
			globalSession: map[string]string{"@message": "Hi John"},
		},
		{
			body:              `set -gF @message "Hi #{@name}"`,
			globalWindowSetup: map[string]string{"@name": "John"},
			globalWindow:      map[string]string{"@name": "John"},
			globalSession:     map[string]string{"@message": "Hi John"},
		},
		{
			body:          `set -gF @message "Hi #{@name}"`,
			sessionSetup:  map[string]string{"@name": "John"},
			session:       map[string]string{"@name": "John"},
			globalSession: map[string]string{"@message": "Hi John"},
		},
		{
			body:               `set -sF @message "Hi #{@name}"`,
			globalSessionSetup: map[string]string{"@name": "John"},
			globalSession:      map[string]string{"@name": "John"},
			server:             map[string]string{"@message": "Hi John"},
		},
		{
			body:          `set -gF @message "Hi #{@name}"`,
			serverSetup:   map[string]string{"@name": "John"},
			server:        map[string]string{"@name": "John"},
			globalSession: map[string]string{"@message": "Hi John"},
		},
	}

	for _, tt := range tests {
		theme := NewTheme()
		s := &SetOptionStatement{}

		if tt.serverSetup != nil {
			theme.ServerOptions = tt.serverSetup
		}
		if tt.globalSessionSetup != nil {
			theme.GlobalSessionOptions = tt.globalSessionSetup
		}
		if tt.sessionSetup != nil {
			theme.SessionOptions = tt.sessionSetup
		}
		if tt.globalWindowSetup != nil {
			theme.GlobalWindowOptions = tt.globalWindowSetup
		}
		if tt.windowSetup != nil {
			theme.WindowOptions = tt.windowSetup
		}

		err := s.Parse(tt.body)
		assert.NoError(t, err)

		err = s.Execute(theme)
		assert.NoError(t, err)

		if tt.server != nil {
			assert.Equal(t, tt.server, theme.ServerOptions)
		}
		if tt.globalSession != nil {
			assert.Equal(t, tt.globalSession, theme.GlobalSessionOptions)
		}
		if tt.session != nil {
			assert.Equal(t, tt.session, theme.SessionOptions)
		}
		if tt.globalWindow != nil {
			assert.Equal(t, tt.globalWindow, theme.GlobalWindowOptions)
		}
		if tt.window != nil {
			assert.Equal(t, tt.window, theme.WindowOptions)
		}
	}
}
