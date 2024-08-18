package resp

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadValue(t *testing.T) {
	t.Run("testing simple string", func(t *testing.T) {
		raw := "+hello\r\n"
		rd := NewReader(bytes.NewBufferString(raw))
		value, err := rd.ReadValue()
		assert.Equal(t, nil, err)
		assert.Equal(t, false, value.IsNull)
		assert.Equal(t, SimpleString, value.Type())
		assert.Equal(t, "hello", string(value.Str))
	})

	t.Run("testing integer", func(t *testing.T) {
		raw := ":1000\r\n"
		rd := NewReader(bytes.NewBufferString(raw))
		value, err := rd.ReadValue()
		assert.Equal(t, nil, err)
		assert.Equal(t, false, value.IsNull)
		assert.Equal(t, Integer, value.Type())
		assert.Equal(t, 1000, value.Integer)
	})

	t.Run("testing bulk string", func(t *testing.T) {
		raw := "$5\r\nhello\r\n"
		rd := NewReader(bytes.NewBufferString(raw))
		value, err := rd.ReadValue()
		assert.Equal(t, nil, err)
		assert.Equal(t, false, value.IsNull)
		assert.Equal(t, BulkString, value.Type())
		assert.Equal(t, "hello", string(value.Str))
	})

	t.Run("testing array", func(t *testing.T) {
		raw := "*3\r\n:1\r\n:2\r\n:3\r\n"
		rd := NewReader(bytes.NewBufferString(raw))
		value, err := rd.ReadValue()
		assert.Equal(t, nil, err)
		assert.Equal(t, false, value.IsNull)
		assert.Equal(t, Array, value.Type())
		assert.Equal(t, 3, len(value.Array))
		for i, value := range value.Array {
			assert.Equal(t, i+1, value.Integer)
		}
	})

	t.Run("testing array", func(t *testing.T) {
		raw := "-errorrrrrrr\r\n"
		rd := NewReader(bytes.NewBufferString(raw))
		value, err := rd.ReadValue()
		assert.Equal(t, nil, err)
		assert.Equal(t, false, value.IsNull)
		assert.Equal(t, Error, value.Type())
		assert.Equal(t, errors.New("errorrrrrrr"), value.Err)
	})

	t.Run("testing wrong resp beginning", func(t *testing.T) {
		raw := "problem\r\n"
		rd := NewReader(bytes.NewBufferString(raw))
		value, err := rd.ReadValue()
		assert.NotEqual(t, nil, err)
		assert.Equal(t, true, value.IsNull)
	})

	t.Run("testing wrong resp ending,no \n", func(t *testing.T) {
		raw := "+problem\r"
		rd := NewReader(bytes.NewBufferString(raw))
		value, err := rd.ReadValue()
		assert.NotEqual(t, nil, err)
		assert.Equal(t, true, value.IsNull)
	})

}
