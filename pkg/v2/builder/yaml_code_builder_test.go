package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYamlCodeBuilder_Bock(t *testing.T) {
	cb := NewYamlCodeBuilder()
	cb.Bock("recipe", func() {
		cb.Line("- banana")
		cb.Line("- bread")
	})
	expected := `recipe:
	- banana
	- bread
`
	assert.Equal(t, expected, cb.String())
}

func TestYamlCodeBuilder_Bockf(t *testing.T) {
	cb := NewYamlCodeBuilder()
	cb.Bockf("%v", "recipe")(func() {
		cb.Line("- banana")
		cb.Line("- bread")
	})
	expected := `recipe:
	- banana
	- bread
`
	assert.Equal(t, expected, cb.String())
}
