package builder

import "fmt"

type GoCodeBuilder interface {
	CodeBuilder
	Block(blockPreamble string, closeWith *string, indentedCode func()) GoCodeBuilder
	DocComment(summary string) GoCodeBuilder
}

type goCodeBuilder struct {
	CodeBuilder
}

func NewGoCodeBuilder() GoCodeBuilder {
	return &goCodeBuilder{
		CodeBuilder: NewCodeBuilder(),
	}
}

// Writes the given text in a nested Go block.
func (g *goCodeBuilder) Block(blockPreamble string, closeWith *string, indentedCode func()) GoCodeBuilder {
	g.CodeBuilder.Append(blockPreamble)
	g.CodeBuilder.Line(" {")
	g.CodeBuilder.IncrementLevel()

	closeWithVal := "}"
	if closeWith != nil {
		closeWithVal += *closeWith
	}
	g.CodeBuilder.WithIndented(&closeWithVal, indentedCode)

	return g
}

// Write a doc comment with the given summary.
func (g *goCodeBuilder) DocComment(summary string) GoCodeBuilder {
	g.BlankLine()
	g.Line(fmt.Sprintf("// %v", summary))
	return g
}
