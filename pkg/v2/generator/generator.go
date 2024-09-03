package generator

import (
	"io"
)

type Generator interface {
	DefaultOutputPath() string
	EffectiveOutputPath() string
	IsDisabled() bool
	Generate() (string, error)
}

type Cleaner interface {
	Clean(knownGoodFiles []string)
}
type ModelGenerator interface {
	Generator
}

type CompositeGenerator interface {
	Generator
	GetGenerators() []Generator
	GetCleaners() []Cleaner
}

type FileGenerator interface {
	Generator
	GetOutput() io.Reader
}

type ModelFileGenerator[TModel any] interface {
	ModelGenerator
	GetOutput() io.Reader
}
