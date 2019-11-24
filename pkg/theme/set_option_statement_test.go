package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatementInterfaceeCompliance(t *testing.T) {
	assert.Implements(t, (*Statement)(nil), &SetOptionStatement{})
}

var setOptionStatementParseTests = []struct {
	body  string
	flags *SetOptionFlags
	args  []string
	error error
}{
	{
		body:  `set -a val1 foo`,
		flags: &SetOptionFlags{Append: true},
		args:  []string{"val1", "foo"},
	},
	{
		body:  `set -F val2 'foo bar'`,
		flags: &SetOptionFlags{Format: true},
		args:  []string{"val2", "foo bar"},
	},
	{
		body:  `set -g val3 "foo bar"`,
		flags: &SetOptionFlags{Global: true},
		args:  []string{"val3", "foo bar"},
	},
	{
		body:  `set -o val4 foo`,
		flags: &SetOptionFlags{OnlyIfUnset: true},
		args:  []string{"val4", "foo"},
	},
	{
		body:  `set -q val5 foo`,
		flags: &SetOptionFlags{Quiet: true},
		args:  []string{"val5", "foo"},
	},
	{
		body:  `set -s val6 foo`,
		flags: &SetOptionFlags{Server: true},
		args:  []string{"val6", "foo"},
	},
	{
		body:  `set -t other:3 val7 foo`,
		flags: &SetOptionFlags{Target: "other:3"},
		args:  []string{"val7", "foo"},
	},
	{
		body:  `set -u val8 foo`,
		flags: &SetOptionFlags{Unset: true},
		args:  []string{"val8", "foo"},
	},
	{
		body:  `set -w val9 foo`,
		flags: &SetOptionFlags{Window: true},
		args:  []string{"val9", "foo"},
	},
	{
		body:  `set-option -w val10 foo`,
		flags: &SetOptionFlags{Window: true},
		args:  []string{"val10", "foo"},
	},
	{
		body:  `set-window-option val11 foo`,
		flags: &SetOptionFlags{Window: true},
		args:  []string{"val11", "foo"},
	},
	{
		body:  `set -goq @val12 'hello world'`,
		flags: &SetOptionFlags{Global: true, OnlyIfUnset: true, Quiet: true},
		args:  []string{"@val12", "hello world"},
	},
	{
		body:  `set -gF @val13 'hello #{@other} world'`,
		flags: &SetOptionFlags{Global: true, Format: true},
		args:  []string{"@val13", "hello #{@other} world"},
	},
	{
		body:  `has-session -t val14`,
		error: &NotSupportedCommandError{"has-session", setOptionCommands},
	},
}

func TestSetOptionStatementParse(t *testing.T) {
	for _, tt := range setOptionStatementParseTests {
		s := &SetOptionStatement{}

		err := s.Parse(tt.body)

		if tt.flags != nil {
			assert.Equal(t, tt.flags, s.Flags)
		}

		if tt.args != nil {
			assert.Equal(t, tt.args, s.Arguments)
		}

		if tt.error != nil {
			assert.Error(t, tt.error, err)
		}
	}
}

var setOptionStatementExecuteTests = []struct {
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
	{
		body:    `set @foo "bar"`,
		session: map[string]string{"@foo": "bar"},
	},
	{
		body:    `set-option @foo "bar"`,
		session: map[string]string{"@foo": "bar"},
	},
	{
		body:   `set -s @foo "bar"`,
		server: map[string]string{"@foo": "bar"},
	},
	{
		body:        `set -so @foo "bar"`,
		serverSetup: map[string]string{"@foo": "foo"},
		server:      map[string]string{"@foo": "foo"},
	},
	{
		body:          `set -g @foo "bar"`,
		globalSession: map[string]string{"@foo": "bar"},
	},
	{
		body:               `set -og @foo "bar"`,
		globalSessionSetup: map[string]string{"@foo": "foo"},
		globalSession:      map[string]string{"@foo": "foo"},
	},
	{
		body:   `set -w @foo "bar"`,
		window: map[string]string{"@foo": "bar"},
	},
	{
		body:        `set -wo @foo "bar"`,
		windowSetup: map[string]string{"@foo": "foo"},
		window:      map[string]string{"@foo": "foo"},
	},
	{
		body:         `set -wg @foo "bar"`,
		globalWindow: map[string]string{"@foo": "bar"},
	},
	{
		body:              `set -wgo @foo "bar"`,
		globalWindowSetup: map[string]string{"@foo": "foo"},
		globalWindow:      map[string]string{"@foo": "foo"},
	},
	{
		body:   `set-window-option @foo "bar"`,
		window: map[string]string{"@foo": "bar"},
	},
	{
		body:        `set-window-option -o @foo "bar"`,
		windowSetup: map[string]string{"@foo": "foo"},
		window:      map[string]string{"@foo": "foo"},
	},
	{
		body:         `set-window-option -g @foo "bar"`,
		globalWindow: map[string]string{"@foo": "bar"},
	},
	{
		body:              `set-window-option -go @foo "bar"`,
		globalWindowSetup: map[string]string{"@foo": "foo"},
		globalWindow:      map[string]string{"@foo": "foo"},
	},
	{
		body:        `set -sa @foo "bar"`,
		serverSetup: map[string]string{"@foo": "foo"},
		server:      map[string]string{"@foo": "foobar"},
	},
	{
		body:               `set -ga @foo "bar"`,
		globalSessionSetup: map[string]string{"@foo": "foo"},
		globalSession:      map[string]string{"@foo": "foobar"},
	},
	{
		body:         `set -a @foo "bar"`,
		sessionSetup: map[string]string{"@foo": "foo"},
		session:      map[string]string{"@foo": "foobar"},
	},
	{
		body:              `set -gwa @foo "bar"`,
		globalWindowSetup: map[string]string{"@foo": "foo"},
		globalWindow:      map[string]string{"@foo": "foobar"},
	},
	{
		body:        `set -wa @foo "bar"`,
		windowSetup: map[string]string{"@foo": "foo"},
		window:      map[string]string{"@foo": "foobar"},
	},
	{
		body:        `set -su @foo`,
		serverSetup: map[string]string{"@foo": "foo", "@bar": "bar"},
		server:      map[string]string{"@bar": "bar"},
	},
	{
		body:               `set -gu @foo`,
		globalSessionSetup: map[string]string{"@foo": "foo", "@bar": "bar"},
		globalSession:      map[string]string{"@bar": "bar"},
	},
	{
		body:         `set -u @foo`,
		sessionSetup: map[string]string{"@foo": "foo", "@bar": "bar"},
		session:      map[string]string{"@bar": "bar"},
	},
	{
		body:              `set -gwu @foo`,
		globalWindowSetup: map[string]string{"@foo": "foo", "@bar": "bar"},
		globalWindow:      map[string]string{"@bar": "bar"},
	},
	{
		body:        `set -wu @foo`,
		windowSetup: map[string]string{"@foo": "foo", "@bar": "bar"},
		window:      map[string]string{"@bar": "bar"},
	},
	{
		body:         `set -u @hello`,
		sessionSetup: map[string]string{"@foo": "foo", "@bar": "bar"},
		session:      map[string]string{"@foo": "foo", "@bar": "bar"},
	},
	{
		body:         `set -F @foo "foo #{@bar} baz"`,
		sessionSetup: map[string]string{"@bar": "bar"},
		session:      map[string]string{"@bar": "bar", "@foo": "foo bar baz"},
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
}

func TestSetOptionStatementExecute(t *testing.T) {
	for _, tt := range setOptionStatementExecuteTests {
		theme := NewTheme()
		s := &SetOptionStatement{theme: theme}

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

		err = s.Execute()
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