package generator

import "io"

type fileGenerator struct {
}

// DefaultOutputPath implements FileGenerator.
func (f *fileGenerator) DefaultOutputPath() string {
	panic("unimplemented")
}

// EffectiveOutputPath implements FileGenerator.
func (f *fileGenerator) EffectiveOutputPath() string {
	panic("unimplemented")
}

// Generate implements FileGenerator.
func (f *fileGenerator) Generate() (string, error) {
	panic("unimplemented")
}

// GetOutput implements FileGenerator.
func (f *fileGenerator) GetOutput() io.Reader {
	panic("unimplemented")
}

// IsDisabled implements FileGenerator.
func (f *fileGenerator) IsDisabled() bool {
	panic("unimplemented")
}

func NewFileGenerator() FileGenerator {
	return &fileGenerator{}
}
