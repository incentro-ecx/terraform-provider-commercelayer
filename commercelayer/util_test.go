package commercelayer

import (
	"fmt"
	"github.com/incentro-dc/go-commercelayer-sdk/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiagErrStandardErr(t *testing.T) {
	msg := "some error"
	diag := diagErr(fmt.Errorf(msg))

	assert.True(t, diag.HasError())
	assert.Equal(t, msg, diag[0].Summary)
}

func TestDiagErrGenericOpenAPIError(t *testing.T) {
	diag := diagErr(&api.GenericOpenAPIError{})
	assert.True(t, diag.HasError())
	assert.Equal(t, ": ", diag[0].Summary)
}

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

func TestKeyValueRefNilVal(t *testing.T) {
	assert.Equal(t, map[string]interface{}{}, keyValueRef(nil))
}

func TestKeyValueRefEmptyMap(t *testing.T) {
	assert.Equal(t, map[string]interface{}{}, keyValueRef(map[string]interface{}{}))
}

func TestKeyValueRefNonEmptyMap(t *testing.T) {
	assert.Equal(t, map[string]interface{}{"hello": "world"}, keyValueRef(map[string]interface{}{"hello": "world"}))
}

func TestBoolRefNilVal(t *testing.T) {
	assert.Nil(t, boolRef(nil))
}

func TestBoolRefTrue(t *testing.T) {
	assert.True(t, *boolRef(true))
}

func TestBoolRefFalse(t *testing.T) {
	assert.False(t, *boolRef(false))
}

func TestFloat64ToFloat32RefNilVal(t *testing.T) {
	assert.Nil(t, float64ToFloat32Ref(nil))
}

func TestFloat64ToFloat32RefZeroVal(t *testing.T) {
	assert.Nil(t, float64ToFloat32Ref(float64(0)))
}

func TestFloat64ToFloat32RefFloat64Val(t *testing.T) {
	assert.Equal(t, float32(1), *float64ToFloat32Ref(float64(1)))
}

func TestStringSliceValueRefNilVal(t *testing.T) {
	assert.Equal(t, []string{}, stringSliceValueRef(nil))
}

func TestStringSliceValueRefEmptySliceVal(t *testing.T) {
	assert.Equal(t, []string(nil), stringSliceValueRef([]interface{}{}))
}

func TestStringSliceValueRefFilledSliceVal(t *testing.T) {
	assert.Equal(t, []string{"foobar"}, stringSliceValueRef([]interface{}{"foobar"}))
}

func TestStringSliceValueRefFilledIntSliceVal(t *testing.T) {
	assert.Panics(t, func() {
		stringSliceValueRef([]interface{}{1})
	})
}
