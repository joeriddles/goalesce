package builder

import (
	"errors"
	"strings"
)

type CodeBuilder interface {
	String() string
	GetLevel() int
	IncrementLevel() CodeBuilder
	DecrementLevel() CodeBuilder
	Line(line string) CodeBuilder
	Lines(lines ...string) CodeBuilder
	BlankLine() CodeBuilder
	Indented(line string) (CodeBuilder, error)
	WithIndented(closeWith *string, indentedCode func()) CodeBuilder
	Append(text string) CodeBuilder
	TrimEnd(text string) CodeBuilder
	TrimWhitespace() CodeBuilder
}

type codeBuilder struct {
	sb         StringBuilder
	level      int
	indentSize int
	indentChar string
	onNewLine  bool
}

func NewCodeBuilder() CodeBuilder {
	return &codeBuilder{
		sb:         NewStringBuilder(),
		level:      0,
		indentSize: 1,
		indentChar: "\t",
		onNewLine:  true,
	}
}

func (c *codeBuilder) GetLevel() int {
	return c.level
}

func (c *codeBuilder) IncrementLevel() CodeBuilder {
	c.level++
	return c
}

func (c *codeBuilder) DecrementLevel() CodeBuilder {
	if c.level > 0 {
		c.level--
	}

	return c
}

// Write a line of text at the current indentation level.
func (c *codeBuilder) Line(line string) CodeBuilder {
	if c.onNewLine {
		c.sb.AppendN(c.indentChar, c.level*c.indentSize)
	}
	c.sb.Append(line).AppendLine()
	c.onNewLine = true
	return c
}

// Calls Line(string) for each line.
func (c *codeBuilder) Lines(lines ...string) CodeBuilder {
	for _, line := range lines {
		c.Line(line)
	}
	return c
}

func (c *codeBuilder) BlankLine() CodeBuilder {
	c.Line("")
	return c
}

// Write nested code that automatically unsets the indentation level.
func (c *codeBuilder) WithIndented(closeWith *string, indentedCode func()) CodeBuilder {
	defer func() {
		c.level--
		if closeWith != nil {
			c.Line(*closeWith)
		}
	}()

	indentedCode()

	return c
}

// Write a line that is indented one level past the current indentation level.
func (c *codeBuilder) Indented(line string) (CodeBuilder, error) {
	if c.onNewLine {
		c.sb.AppendN(c.indentChar, (c.level+1)*c.indentSize)
	} else {
		return c, errors.New("cannot start an indented line on a line that isn't empty")
	}
	c.sb.Append(line).AppendLine()
	c.onNewLine = true
	return c, nil
}

// Write text to the current line. If currently on a new, blank line, the current indentation will be added.
func (c *codeBuilder) Append(text string) CodeBuilder {
	if c.onNewLine {
		c.sb.AppendN(c.indentChar, c.level*c.indentSize)
		c.onNewLine = false
	}
	c.sb.Append(text)
	return c
}

// Trim the given string from the end of the output if it exists.
func (c *codeBuilder) TrimEnd(text string) CodeBuilder {
	str := c.sb.String()
	start := c.sb.GetLength() - len(text)
	if start < 0 {
		return c
	}
	str = str[start : start+len(text)]
	if str == text {
		c.sb.Remove(start, len(text))
	}
	return c
}

var asciiSpace = " \t\n\v\f\r" // taken from strings.go

// Trim whitespace from the end of the output.
func (c *codeBuilder) TrimWhitespace() CodeBuilder {
	sb := NewStringBuilder()
	sb.Append(strings.TrimRight(c.sb.String(), asciiSpace))
	c.sb = sb
	return c
}

func (c *codeBuilder) String() string {
	return c.sb.String()
}
