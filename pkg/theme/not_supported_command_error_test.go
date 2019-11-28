package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotSupportedCommandErrorInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*error)(nil), &NotSupportedCommandError{})
}

func TestNotSupportedCommandError(t *testing.T) {
	var tests = []struct {
		command           string
		supportedCommands []string
		msg               string
	}{
		{
			"foo", []string{"bar", "baz"},
			"foo is not one of the supported commands: bar, baz",
		},
		{
			"has-session", []string{"set", "set-option", "set-window-option"},
			"has-session is not one of the supported commands: " +
				"set, set-option, set-window-option",
		},
	}

	for _, tt := range tests {
		err := NotSupportedCommandError{
			Command:           tt.command,
			SupportedCommands: tt.supportedCommands,
		}

		assert.Equal(t, tt.msg, err.Error())
	}
}
