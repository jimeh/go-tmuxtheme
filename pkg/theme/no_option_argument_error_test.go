package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoOptionArgumentErrorInterfaceCompliance(t *testing.T) {
	assert.Implements(t, (*error)(nil), &NoOptionArgumentError{})
}

func TestNoOptionArgumentError(t *testing.T) {
	err := &NoOptionArgumentError{}

	assert.Equal(t, "No option argument given", err.Error())
}
