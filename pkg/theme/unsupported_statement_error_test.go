package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsupportedStatementErrorInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*error)(nil), &UnsupportedStatementError{})
}

var unsupportedStatementErrorTests = []struct {
	body string
	err  string
}{
	{"foo", "Unsupported statement: foo"},
	{"has-session -t other:3", "Unsupported statement: has-session -t other:3"},
}

func TestUnsupportedStatementError(t *testing.T) {
	for _, tt := range unsupportedStatementErrorTests {
		err := UnsupportedStatementError{Body: tt.body}

		assert.Equal(t, tt.err, err.Error())
	}
}
