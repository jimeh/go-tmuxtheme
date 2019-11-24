package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyStatementInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*Statement)(nil), &EmptyStatement{})
}

var emptyStatementParseTests = []struct {
	body  string
	error error
}{
	{body: ``},
	{
		body:  `# This is a comment`,
		error: &NotSupportedCommandError{"", []string{}},
	},
	{
		body:  `set -g foo "bar"`,
		error: &NotSupportedCommandError{"", []string{}},
	},
}

func TestEmptyStatementParse(t *testing.T) {
	for _, tt := range emptyStatementParseTests {
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
