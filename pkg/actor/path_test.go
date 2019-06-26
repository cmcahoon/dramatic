package actor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPath_String(t *testing.T) {
	path := &Path{
		"bar",
		"foo",
		"",
	}
	pathString := path.String()

	assert.Equal(t, "/foo/bar", pathString)
}
