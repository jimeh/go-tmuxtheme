package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyStatementInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*Statement)(nil), &EmptyStatement{})
}

func TestEmptyStatementParse(t *testing.T) {
	var tests = []struct {
		body  string
		error error
	}{
		{body: ``},
		{
			body:  `# This is a comment`,
			error: &NotSupportedCommandError{"#", []string{}},
		},
		{
			body:  `# it's a comment`,
			error: &NotSupportedCommandError{"#", []string{}},
		},
		{
			body:  `  # it's a comment`,
			error: &NotSupportedCommandError{"#", []string{}},
		},
		{
			body:  `set -g foo "bar"`,
			error: &NotSupportedCommandError{"set", []string{}},
		},
	}

	for _, tt := range tests {
		s := &EmptyStatement{}

		err := s.Parse(tt.body)

		if tt.error != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.error, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
