package builder

type YamlCodeBuilder interface {
	CodeBuilder
	Bock(blockPreamble string, indentedCode func()) YamlCodeBuilder
	Bockf(blockPreamble string, args ...any) func(func())
	List(blockPreamble string, indentedCode func()) YamlCodeBuilder
	DocComment(summary string) YamlCodeBuilder
}

type yamlCodeBuilder struct {
	CodeBuilder
	thenCount int
}

func NewYamlCodeBuilder() YamlCodeBuilder {
	return &yamlCodeBuilder{
		CodeBuilder: NewCodeBuilder(2, " "),
		thenCount:   0,
	}
}

func (y *yamlCodeBuilder) Bock(blockPreamble string, indentedCode func()) YamlCodeBuilder {
	y.CodeBuilder.Line(blockPreamble + ":")
	y.CodeBuilder.IncrementLevel()
	y.CodeBuilder.WithIndented(nil, indentedCode)
	return y
}

func (y *yamlCodeBuilder) Bockf(blockPreamble string, args ...any) func(func()) {
	y.CodeBuilder.Linef(blockPreamble+":", args...)
	y.CodeBuilder.IncrementLevel()
	return func(indentedCode func()) {
		y.CodeBuilder.WithIndented(nil, indentedCode)
	}
}

func (y *yamlCodeBuilder) List(blockPreamble string, indentedCode func()) YamlCodeBuilder {
	y.CodeBuilder.Line("- " + blockPreamble)
	y.CodeBuilder.IncrementLevel()
	y.CodeBuilder.WithIndented(nil, indentedCode)
	return y
}

func (y *yamlCodeBuilder) DocComment(summary string) YamlCodeBuilder {
	y.CodeBuilder.Linef("# %v", summary)
	return y
}
