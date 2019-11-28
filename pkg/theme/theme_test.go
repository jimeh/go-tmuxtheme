package theme

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThemeParse(t *testing.T) {
	var tests = []struct {
		body       string
		result     []string
		statements []Statement
	}{
		{
			body: `
set -g @name "John Smith"
set -gF @message "Hi #{@name}"
`,
			statements: []Statement{
				&SetOptionStatement{
					Option: "@name",
					Value:  "John Smith",
					Flags:  &SetOptionFlags{Global: true},
				},
				&SetOptionStatement{
					Option: "@message",
					Value:  "Hi #{@name}",
					Flags:  &SetOptionFlags{Global: true, Format: true},
				},
			},
		},
		{
			body: `
set -g @name "John Smith"
set -gF @message \
  "Hi #{@name}"

`,
			statements: []Statement{
				&SetOptionStatement{
					Option: "@name",
					Value:  "John Smith",
					Flags:  &SetOptionFlags{Global: true},
				},
				&SetOptionStatement{
					Option: "@message",
					Value:  "Hi #{@name}",
					Flags:  &SetOptionFlags{Global: true, Format: true},
				},
				&EmptyStatement{},
			},
		},
		{
			body: `
set -g @name "John Smith"

# This is the message
set -gF @message \
  "Hi #{@name}"
`,
			statements: []Statement{
				&SetOptionStatement{
					Option: "@name",
					Value:  "John Smith",
					Flags:  &SetOptionFlags{Global: true},
				},
				&EmptyStatement{},
				&CommentStatement{
					Msg: "This is the message",
				},
				&SetOptionStatement{
					Option: "@message",
					Value:  "Hi #{@name}",
					Flags:  &SetOptionFlags{Global: true, Format: true},
				},
			},
		},
	}

	for _, tt := range tests {
		theme := New()
		r := strings.NewReader(tt.body[1:])

		err := theme.Parse(r)

		assert.NoError(t, err)
		assert.Equal(t, tt.statements, theme.Statements)
	}
}

func TestThemeExecute(t *testing.T) {
	var tests = []struct {
		body          string
		server        map[string]string
		globalSession map[string]string
		session       map[string]string
		globalWindow  map[string]string
		window        map[string]string
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
			body: `
set @name "John"
set -o @name "Jim"
`,
			session: map[string]string{"@name": "John"},
		},
		{
			body: `
set @name "John"
set -a @name "Jim"
`,
			session: map[string]string{"@name": "JohnJim"},
		},
		{
			body: `
set @name "John"
set @foo "bar"
set -u @name
`,
			session: map[string]string{"@foo": "bar"},
		},
		{
			body: `
set @name 'John Smith'
set -F @message "Hi #{@name}"
`,
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
			body: `
set -s @name "John"
set -so @name "Jim"
`,
			server: map[string]string{"@name": "John"},
		},
		{
			body: `
set -s @name "John"
set -sa @name "Jim"
`,
			server: map[string]string{"@name": "JohnJim"},
		},
		{
			body: `
set -s @name "John"
set -s @foo "bar"
set -su @name
`,
			server: map[string]string{"@foo": "bar"},
		},
		{
			body: `
set -s @name 'John Smith'
set -sF @message "Hi #{@name}"
`,
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
			body: `
set -g @name "John"
set -go @name "Jim"
`,
			globalSession: map[string]string{"@name": "John"},
		},
		{
			body: `
set -g @name "John"
set -ga @name "Jim"
`,
			globalSession: map[string]string{"@name": "JohnJim"},
		},
		{
			body: `
set -g @name "John"
set -gu @name
`,
			globalSession: map[string]string{},
		},
		{
			body: `
set -g @name 'John Smith'
set -gF @message "Hi #{@name}"
`,
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
			body: `
set -w @name "John"
set -wo @name "Jim"
`,
			window: map[string]string{"@name": "John"},
		},
		{
			body: `
set -w @name "John"
set -wa @name "Jim"
`,
			window: map[string]string{"@name": "JohnJim"},
		},
		{
			body: `
set -w @name "John"
set -w @foo "bar"
set -wu @name
`,
			window: map[string]string{"@foo": "bar"},
		},
		{
			body: `
set -w @name 'John Smith'
set -wF @message "Hi #{@name}"
`,
			window: map[string]string{
				"@name":    "John Smith",
				"@message": "Hi John Smith",
			},
		},
		//
		// Global Window Options
		//
		{
			body:         `set -wg @name "John"`,
			globalWindow: map[string]string{"@name": "John"},
		},
		{
			body:         `set-option -gw @name "John"`,
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
			body: `
set -gw @name "John"
set -gwo @name "Jim"
`,
			globalWindow: map[string]string{"@name": "John"},
		},
		{
			body: `
set -wg @name "John"
set -wga @name "Jim"
`,
			globalWindow: map[string]string{"@name": "JohnJim"},
		},
		{
			body: `
set -gw @name "John"
set -gw @foo "bar"
set -gwu @name
`,
			globalWindow: map[string]string{"@foo": "bar"},
		},
		{
			body: `
set -wg @name 'John Smith'
set -wgF @message "Hi #{@name}"
`,
			globalWindow: map[string]string{
				"@name":    "John Smith",
				"@message": "Hi John Smith",
			},
		},
		//
		// Formatting
		//
		{
			body: `
set @bar "bar"
set -F @foo "foo #{@bar} baz"
`,
			session: map[string]string{"@bar": "bar", "@foo": "foo bar baz"},
		},
		{
			body: `
set bar 'bar'
set -F @foo "foo #{bar} baz"
`,
			session: map[string]string{"bar": "bar", "@foo": "foo bar baz"},
		},
		{
			body: `
set @bar bar
set -F @foo "foo #{@bar}#{@bar}"
`,
			session: map[string]string{"@bar": "bar", "@foo": "foo barbar"},
		},
		{
			body: `
set @foo foo
set @bar bar
set -F @msg "#{@foo} #{@bar}"
`,
			session: map[string]string{
				"@foo": "foo",
				"@bar": "bar",
				"@msg": "foo bar",
			},
		},
		{
			body: `
set -g @bar "bar"
set -gF @foo "foo #{@bar} baz"
`,
			globalSession: map[string]string{
				"@bar": "bar",
				"@foo": "foo bar baz",
			},
		},
		{
			body: `
set @bar bar
set @foo foo
set -Fo @foo "foo #{@bar} baz"
`,
			session: map[string]string{"@bar": "bar", "@foo": "foo"},
		},
		{
			body: `
set @bar bar
set @foo foo
set -Fa @foo " #{@bar} baz"
`,
			session: map[string]string{"@bar": "bar", "@foo": "foo bar baz"},
		},
		{
			body: `
set -w @name "John"
set -gF @message "Hi #{@name}"
`,
			window:        map[string]string{"@name": "John"},
			globalSession: map[string]string{"@message": "Hi John"},
		},
		{
			body: `
set -gw @name John
set -gF @message "Hi #{@name}"
`,
			globalWindow:  map[string]string{"@name": "John"},
			globalSession: map[string]string{"@message": "Hi John"},
		},
		{
			body: `
set @name 'John'
set -gF @message "Hi #{@name}"
`,
			session:       map[string]string{"@name": "John"},
			globalSession: map[string]string{"@message": "Hi John"},
		},
		{
			body: `
set -g @name "John"
set -sF @message "Hi #{@name}"
`,
			globalSession: map[string]string{"@name": "John"},
			server:        map[string]string{"@message": "Hi John"},
		},
		{
			body: `
set -s @name John
set -gF @message "Hi #{@name}"
`,
			server:        map[string]string{"@name": "John"},
			globalSession: map[string]string{"@message": "Hi John"},
		},
	}

	for _, tt := range tests {
		theme := New()
		r := strings.NewReader(strings.TrimLeft(tt.body, "\n"))

		err := theme.Parse(r)
		assert.NoError(t, err)

		err = theme.Execute()
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

func TestThemeLoadFileAndExecute(t *testing.T) {
	theme := New()

	err := theme.Load("theme_test.tmuxtheme")
	require.NoError(t, err)

	err = theme.Execute()
	require.NoError(t, err)

	assert.Equal(t, theme.ServerOptions, map[string]string{})
	assert.Equal(
		t,
		theme.GlobalSessionOptions,
		map[string]string{
			"@theme-clock-mode-colour":                   "red",
			"@theme-clock-mode-style":                    "24",
			"@theme-display-panes-active-colour":         "default",
			"@theme-display-panes-colour":                "default",
			"@theme-message-bg":                          "default",
			"@theme-message-command-bg":                  "default",
			"@theme-message-command-fg":                  "default",
			"@theme-message-fg":                          "default",
			"@theme-mode-bg":                             "red",
			"@theme-mode-fg":                             "default",
			"@theme-pane-active-border-bg":               "default",
			"@theme-pane-active-border-fg":               "green",
			"@theme-pane-border-bg":                      "default",
			"@theme-pane-border-fg":                      "default",
			"@theme-status-bg":                           "black",
			"@theme-status-fg":                           "cyan",
			"@theme-status-interval":                     "1",
			"@theme-status-justify":                      "centre",
			"@theme-status-left":                         "#S #[fg=white]» #[fg=yellow]#I #[fg=cyan]#P",
			"@theme-status-left-bg":                      "black",
			"@theme-status-left-fg":                      "green",
			"@theme-status-left-length":                  "40",
			"@theme-status-right":                        "#H #[fg=white]« #[fg=yellow]%H:%M:%S #[fg=green]%d-%b-%y",
			"@theme-status-right-bg":                     "black",
			"@theme-status-right-fg":                     "cyan",
			"@theme-status-right-length":                 "40",
			"@theme-window-status-activity-bg":           "black",
			"@theme-window-status-activity-fg":           "yellow",
			"@theme-window-status-current-bg":            "red",
			"@theme-window-status-current-fg":            "black",
			"@theme-window-status-current-format":        " #I:#W#F ",
			"@theme-window-status-format":                " #I:#W#F ",
			"@theme-window-status-separator":             "",
			"@themepack-status-left-area-left-format":    "#S",
			"@themepack-status-left-area-middle-format":  "#I",
			"@themepack-status-left-area-right-format":   "#P",
			"@themepack-status-right-area-left-format":   "#H",
			"@themepack-status-right-area-middle-format": "%H:%M:%S",
			"@themepack-status-right-area-right-format":  "%d-%b-%y",
			"@themepack-window-status-current-format":    "#I:#W#F",
			"@themepack-window-status-format":            "#I:#W#F",
			"clock-mode-colour":                          "red",
			"clock-mode-style":                           "24",
			"display-panes-active-colour":                "default",
			"display-panes-colour":                       "default",
			"message-command-style":                      "bg=default,fg=default",
			"message-style":                              "bg=default,fg=default",
			"mode-style":                                 "bg=red,fg=default",
			"pane-active-border-style":                   "bg=default,fg=green",
			"pane-border-style":                          "bg=default,fg=default",
			"status-interval":                            "1",
			"status-justify":                             "centre",
			"status-left":                                "#S #[fg=white]» #[fg=yellow]#I #[fg=cyan]#P",
			"status-left-length":                         "40",
			"status-left-style":                          "bg=black,fg=green",
			"status-right":                               "#H #[fg=white]« #[fg=yellow]%H:%M:%S #[fg=green]%d-%b-%y",
			"status-right-length":                        "40",
			"status-right-style":                         "bg=black,fg=cyan",
			"status-style":                               "bg=black,fg=cyan",
			"window-status-activity-style":               "bg=black,fg=yellow",
			"window-status-current-format":               " #I:#W#F ",
			"window-status-current-style":                "bg=red,fg=black",
			"window-status-format":                       " #I:#W#F ",
			"window-status-separator":                    "",
		},
	)
	assert.Equal(t, theme.SessionOptions, map[string]string{})
	assert.Equal(t, theme.GlobalWindowOptions, map[string]string{})
	assert.Equal(t, theme.WindowOptions, map[string]string{})
}
