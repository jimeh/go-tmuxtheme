package theme

import "testing"

import "github.com/stretchr/testify/assert"

var newStatementTests = []struct {
	body      string
	statement Statement
	error     error
}{
	// SetOptionStatement
	{
		body: `set foo bar`,
		statement: &SetOptionStatement{
			Flags:  &SetOptionFlags{},
			Option: "foo",
			Value:  "bar",
		},
	},
	{
		body: `set-option foo bar`,
		statement: &SetOptionStatement{
			Flags:  &SetOptionFlags{},
			Option: "foo",
			Value:  "bar",
		},
	},
	{
		body: `  set-window-option foo bar`,
		statement: &SetOptionStatement{
			Flags:  &SetOptionFlags{Window: true},
			Option: "foo",
			Value:  "bar",
		},
	},
	// CommentStatement
	{
		body:      `# This is a comment`,
		statement: &CommentStatement{Msg: "This is a comment"},
	},
	{
		body:      `#This is a comment`,
		statement: &CommentStatement{Msg: "This is a comment"},
	},
	{
		body:      `  # This is a comment`,
		statement: &CommentStatement{Msg: "This is a comment"},
	},
	// EmptyStatement
	{
		body:      ``,
		statement: &EmptyStatement{},
	},
	{
		body:      `  `,
		statement: &EmptyStatement{},
	},
	// Unsupported Statement
	{
		body:  `has-session -t other:3`,
		error: &UnsupportedStatementError{Body: "has-session -t other:3"},
	},
}

func TestNewStatement(t *testing.T) {
	for _, tt := range newStatementTests {
		st, err := NewStatement(tt.body)

		if tt.statement != nil {
			assert.NoError(t, err)
			assert.Equal(t, tt.statement, st)
		}

		if tt.error != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.error, err)
		}
	}
}
