package commercelayer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringRefNilVal(t *testing.T) {
	assert.Nil(t, stringRef(nil))
}

func TestStringRefEmptyStringVal(t *testing.T) {
	assert.Nil(t, stringRef(""))
}

func TestStringRefExistingStringVal(t *testing.T) {
	str := "foobar"
	strRef := stringRef(str)
	assert.NotNil(t, strRef)
	assert.Equal(t, strRef, &str)
}
