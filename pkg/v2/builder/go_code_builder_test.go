package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Block(t *testing.T) {
	cb := NewGoCodeBuilder()
	cb.Block("if 1 == 1", nil, func() {
		cb.Line(`fmt.Println("hello world")`)
	})
	expected := `if 1 == 1 {
	fmt.Println("hello world")
}
`
	assert.Equal(t, expected, cb.String())
}

func Test_Block_WithCloseWith(t *testing.T) {
	cb := NewGoCodeBuilder()
	closeWith := ` // cool comment`
	cb.Block("if 1 == 1", &closeWith, func() {
		cb.Line(`fmt.Println("hello world")`)
	})
	expected := `if 1 == 1 {
	fmt.Println("hello world")
} // cool comment
`
	assert.Equal(t, expected, cb.String())
}

func Test_DocComment(t *testing.T) {
	cb := NewGoCodeBuilder()
	cb.DocComment("this is a really great function")
	assert.Equal(t, "\n// this is a really great function\n", cb.String())
}
