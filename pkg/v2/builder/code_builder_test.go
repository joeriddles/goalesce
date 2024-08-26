package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Level(t *testing.T) {
	cb := NewCodeBuilder()
	assert.Equal(t, 0, cb.GetLevel())
	assert.Equal(t, 1, cb.IncrementLevel().GetLevel())
	assert.Equal(t, 0, cb.DecrementLevel().GetLevel())
	assert.Equal(t, 0, cb.DecrementLevel().GetLevel())
}

func Test_Line(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Line("hello world")
	assert.Equal(t, "hello world\n", cb.String())
	cb.Line("goodbye")
	assert.Equal(t, "hello world\ngoodbye\n", cb.String())
}

func Test_Lines(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Lines("hello world", "goodbye")
	assert.Equal(t, "hello world\ngoodbye\n", cb.String())
}

func Test_BlankLine(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Line("hello")
	cb.BlankLine()
	cb.Line("world")
	assert.Equal(t, "hello\n\nworld\n", cb.String())
}

func Test_Indented_OnNewLine(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Line("if {")
	cb.Indented("hello world")
	assert.Equal(t, "if {\n\thello world\n", cb.String())
}

func Test_Indented_NotOnNewLine(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("not empty")
	_, err := cb.Indented("hello world")
	require.Error(t, err)
}

func Test_WithIndented_NoCloseWith(t *testing.T) {
	cb := NewCodeBuilder()
	cb.WithIndented(nil, func() {
		cb.Line("if {")
		cb.Indented("hello world")
	})
	expected := `if {
	hello world
`
	assert.Equal(t, expected, cb.String())
}

func Test_WithIndented_CloseWith(t *testing.T) {
	cb := NewCodeBuilder()
	closeWith := "}"
	cb.WithIndented(&closeWith, func() {
		cb.Line("if {")
		cb.Indented("hello world")
	})
	expected := `if {
	hello world
}
`
	assert.Equal(t, expected, cb.String())
}

func Test_Append_OnNewLine(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("hello world")
	assert.Equal(t, "hello world", cb.String())
}

func Test_Append_IndentedNewLine(t *testing.T) {
	cb := NewCodeBuilder()
	cb.IncrementLevel()
	cb.Append("hello world")
	assert.Equal(t, "\thello world", cb.String())
}

func Test_Append_NotOnNewLine(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("hello")
	cb.Append(" world")
	assert.Equal(t, "hello world", cb.String())
}

func Test_TrimEnd_Match(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("hello world")
	cb.TrimEnd(" world")
	assert.Equal(t, "hello", cb.String())
}

func Test_TrimEnd_NoMatch(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("hello world")
	cb.TrimEnd(" nope")
	assert.Equal(t, "hello world", cb.String())
}

func Test_TrimEnd_WholeString(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("hello world")
	cb.TrimEnd("hello world")
	assert.Equal(t, "", cb.String())
}

func Test_TrimEnd_TooFar(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("hello world")
	cb.TrimEnd("abc hello world")
	assert.Equal(t, "hello world", cb.String())
}

func Test_TrimeWhitespace(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append("hello world ")
	cb.TrimWhitespace()
	assert.Equal(t, "hello world", cb.String())
}

func Test_TrimeWhitespace_All(t *testing.T) {
	cb := NewCodeBuilder()
	cb.Append(`hello world

    `)
	cb.TrimWhitespace()
	assert.Equal(t, "hello world", cb.String())
}
