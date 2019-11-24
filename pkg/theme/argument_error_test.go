package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentErrorInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*error)(nil), &ArgumentError{})
}

var argumentErrorTests = []struct {
	msg string
}{
	{"No option argument given"},
	{"Something foo bar baz"},
}

func TestArgumentError(t *testing.T) {
	for _, tt := range argumentErrorTests {
		err := ArgumentError{tt.msg}

		assert.Equal(t, tt.msg, err.Error())
	}
}
