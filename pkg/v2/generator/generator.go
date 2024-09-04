package generator

import (
	"io"
)

type Generator[T any] interface {
	DefaultOutputPath() string
	EffectiveOutputPath() string
	IsDisabled() bool
	Generate(T) (string, error)
}

type Cleaner interface {
	Clean(knownGoodFiles []string)
}
type ModelGenerator[T any] interface {
	Generator[T]
}

type CompositeGenerator[T any] interface {
	Generator[T]
	GetGenerators() []Generator[T]
	GetCleaners() []Cleaner
}

type FileGenerator[T any] interface {
	Generator[T]
	GetOutput() io.Reader
}

type ModelFileGenerator[T any] interface {
	ModelGenerator[T]
	GetOutput() io.Reader
}
