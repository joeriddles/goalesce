package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Append(t *testing.T) {
	sb := NewStringBuilder()
	actual := sb.Append("hello").Append("world").String()
	assert.Equal(t, "helloworld", actual)
}

func Test_AppendLine(t *testing.T) {
	sb := NewStringBuilder()
	actual := sb.Append("hello").AppendLine().Append("world").String()
	assert.Equal(t, "hello\nworld", actual)
}

func Test_AppendN(t *testing.T) {
	sb := NewStringBuilder()
	actual := sb.AppendN("hello world. ", 5).String()
	assert.Equal(t, "hello world. hello world. hello world. hello world. hello world. ", actual)
}

func Test_GetLength(t *testing.T) {
	sb := NewStringBuilder()
	assert.Equal(t, 0, sb.GetLength())

	sb.Append("hello")
	assert.Equal(t, 5, sb.GetLength())

	sb.AppendLine()
	assert.Equal(t, 6, sb.GetLength())
}

func Test_Remove(t *testing.T) {
	sb := NewStringBuilder()
	sb.Append("hello world")

	sb.Remove(0, 6)
	assert.Equal(t, "world", sb.String())

	sb.Remove(2, -1)
	assert.Equal(t, "woorld", sb.String())
}
