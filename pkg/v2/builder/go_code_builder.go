package builder

import (
	"fmt"

	"golang.org/x/tools/imports"
)

type GoCodeBuilder interface {
	CodeBuilder
	Bock(blockPreamble string, indentedCode func()) GoCodeBuilder
	Bockf(blockPreamble string, args ...any) func(func())
	ImportBlock(imports ...string) GoCodeBuilder
	DocComment(summary string) GoCodeBuilder
	ErrCheck(ret string) CodeBuilder
	ErrCheckNil() CodeBuilder
	Format() (string, error)
}

type goCodeBuilder struct {
	CodeBuilder
}

func NewGoCodeBuilder() GoCodeBuilder {
	return &goCodeBuilder{
		CodeBuilder: NewCodeBuilder(1, "\t"),
	}
}

// Writes the given text in a nested Go block.
// Purposefully name `Bock` instead of `Block` because it lines up better with
// the four-character function `Line`.
func (g *goCodeBuilder) Bock(blockPreamble string, indentedCode func()) GoCodeBuilder {
	g.CodeBuilder.Append(blockPreamble)
	g.CodeBuilder.Line(" {")
	g.CodeBuilder.IncrementLevel()

	closeWith := "}"
	g.CodeBuilder.WithIndented(&closeWith, indentedCode)

	return g
}

// Writes the formatted text in a nested Go block.
// Returns a callable to write code inside the block.
func (g *goCodeBuilder) Bockf(blockPreamble string, args ...any) func(func()) {
	g.CodeBuilder.Append(fmt.Sprintf(blockPreamble, args...))
	g.CodeBuilder.Line(" {")
	g.CodeBuilder.IncrementLevel()

	return func(indentedCode func()) {
		closeWith := "}"
		g.CodeBuilder.WithIndented(&closeWith, indentedCode)
	}
}

// Writes the given imports.
func (g *goCodeBuilder) ImportBlock(imports ...string) GoCodeBuilder {
	g.CodeBuilder.Line("import (")
	g.CodeBuilder.IncrementLevel()

	closeWithVal := ")"
	g.CodeBuilder.WithIndented(&closeWithVal, func() {
		for _, imprt := range imports {
			g.CodeBuilder.Line(imprt)
		}
	})

	return g
}

// Write a doc comment with the given summary.
func (g *goCodeBuilder) DocComment(summary string) GoCodeBuilder {
	g.Line(fmt.Sprintf("// %v", summary))
	return g
}

// Check and return if `err` is not nil.
func (g *goCodeBuilder) ErrCheck(ret string) CodeBuilder {
	g.Bock("if err != nil", func() {
		g.Linef("return %v, err", ret)
	})
	return g
}

// Check and return nil if `err` is not nil.
func (g *goCodeBuilder) ErrCheckNil() CodeBuilder {
	g.Bock("if err != nil", func() {
		g.Line("return nil, err")
	})
	return g
}

func (g *goCodeBuilder) Format() (string, error) {
	code := g.String()
	bytes := []byte(code)
	formattedBytes, err := imports.Process("", bytes, nil)
	formattedCode := ""
	if err == nil {
		formattedCode = string(formattedBytes)
	}
	return formattedCode, err
}
