package builder

import (
	"errors"
	"strings"
)

type CodeBuilder interface {
	GetLevel() int
	Line(line string) CodeBuilder
	Lines(lines ...string) CodeBuilder
	Indented(line string) (CodeBuilder, error)
}

type codeBuilder struct {
	sb         StringBuilder
	level      int
	indentSize int
	indentChar string
	onNewLine  bool
}

var _ CodeBuilder = new(codeBuilder)

func NewCodeBuilder() CodeBuilder {
	return &codeBuilder{}
}

func (c *codeBuilder) GetLevel() int {
	return c.level
}

// Write a line of text at the current indentation level.
func (c *codeBuilder) Line(line string) CodeBuilder {
	if c.onNewLine {
		c.sb.Append(c.indentChar, c.level*c.indentSize)
	}
	c.sb.Append(line, 1).AppendLine()
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

// Write a line that is indented one level past the current indentation level.
func (c *codeBuilder) Indented(line string) (CodeBuilder, error) {
	if c.onNewLine {
		c.sb.Append(c.indentChar, (c.level+1)*c.indentSize)
	} else {
		return c, errors.New("cannot start an indented line on a line that isn't empty")
	}
	c.sb.Append(line, 1).AppendLine()
	c.onNewLine = true
	return c, nil
}

// Write text to the current line. If currently on a new, blank line, the current indentation will be added.
func (c *codeBuilder) Append(text string) CodeBuilder {
	if c.onNewLine {
		c.sb.Append(c.indentChar, c.level*c.indentSize)
	}
	c.sb.Append(text, 1)
	return c
}

// Trim the given string from the end of the output if it exists.
func (c *codeBuilder) TimeEnd(text string) CodeBuilder {
	str := c.sb.String()
	start := c.sb.GetLength() - len(text)
	end := len(text)
	str = str[start:end]
	if str == text {
		c.sb.Remove(start, end)
	}
	return c
}

// Trim whitespace from the end of the output.
func (c *codeBuilder) TrimWhitespace() CodeBuilder {
	count := 0
	for {
		str := c.sb.String()[c.sb.GetLength()-1-count:]
		if c.sb.GetLength()-count > 0 && c.isWhitespace(str) {
			count++
		} else { // TODO: refactor to use Go syntax for a while loop
			break
		}
	}

	if count > 0 {
		c.sb.Remove(c.sb.GetLength()-count, count)
	}
	return c
}

func (c *codeBuilder) String() string {
	return c.sb.String()
}

func (c *codeBuilder) isWhitespace(text string) bool {
	return len(strings.TrimSpace(text)) == 0
}
