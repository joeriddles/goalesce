package generator

import "io"

type fileGenerator[T any] struct {
}

// DefaultOutputPath implements FileGenerator.
func (f *fileGenerator[T]) DefaultOutputPath() string {
	panic("unimplemented")
}

// EffectiveOutputPath implements FileGenerator.
func (f *fileGenerator[T]) EffectiveOutputPath() string {
	panic("unimplemented")
}

// Generate implements FileGenerator.
func (f *fileGenerator[T]) Generate(_ T) (string, error) {
	panic("unimplemented")
}

// GetOutput implements FileGenerator.
func (f *fileGenerator[T]) GetOutput() io.Reader {
	panic("unimplemented")
}

// IsDisabled implements FileGenerator.
func (f *fileGenerator[T]) IsDisabled() bool {
	panic("unimplemented")
}

func NewFileGenerator[T any]() FileGenerator[T] {
	return &fileGenerator[T]{}
}
