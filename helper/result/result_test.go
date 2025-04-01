package result_test

import (
	"errors"
	"testing"

	"github.com/lazyfury/pulse/helper/result"
	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	assert := assert.New(t)

	// ok must no err
	val := result.Ok(1)
	assert.Equal(1, val.Value)
	assert.Nil(val.Error)

	// Err str
	val2 := result.Err[any](errors.New("err"))
	assert.Equal("err", val2.Error.Error())

	//Err must err
	val3 := result.Err[any](nil)
	assert.Equal(true, val3.IsErr())
}
