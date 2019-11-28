package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentErrorInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*error)(nil), &ArgumentError{})
}

func TestArgumentError(t *testing.T) {
	var tests = []struct {
		msg string
	}{
		{"No option argument given"},
		{"Something foo bar baz"},
	}

	for _, tt := range tests {
		err := ArgumentError{tt.msg}

		assert.Equal(t, tt.msg, err.Error())
	}
}
