package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentStatementInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*Statement)(nil), &CommentStatement{})
}

func TestCommentStatementParse(t *testing.T) {
	var tests = []struct {
		body  string
		msg   string
		error error
	}{
		{body: `# This is a comment`, msg: "This is a comment"},
		{body: `#  This is a comment`, msg: "This is a comment"},
		{body: `#  This is a comment `, msg: "This is a comment"},
		{body: `#This is a comment`, msg: "This is a comment"},
		{body: `#This is a comment `, msg: "This is a comment"},
		{body: `#This is a comment  `, msg: "This is a comment"},
		{body: `  # This is a comment`, msg: "This is a comment"},
		{body: `  #  This is a comment`, msg: "This is a comment"},
		{body: `  #  This is a comment `, msg: "This is a comment"},
		{body: `  #This is a comment`, msg: "This is a comment"},
		{body: `  #This is a comment `, msg: "This is a comment"},
		{body: `  #This is a comment  `, msg: "This is a comment"},
		{body: `#`, msg: ""},
		{body: `   #`, msg: ""},
		{body: `#    `, msg: ""},
		{
			body:  ``,
			error: &NotSupportedCommandError{"", []string{"#"}},
		},
		{
			body:  `set -g @foo "bar"`,
			error: &NotSupportedCommandError{"set", []string{"#"}},
		},
		{
			body:  `set -g @foo "bar" # This is a comment`,
			error: &NotSupportedCommandError{"set", []string{"#"}},
		},
	}

	for _, tt := range tests {
		s := &CommentStatement{}

		err := s.Parse(tt.body)

		if tt.msg != "" {
			assert.Equal(t, tt.msg, s.Msg)
		}

		if tt.error != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.error, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
