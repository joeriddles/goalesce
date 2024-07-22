package config

import (
	"os"
	"path/filepath"
)

type Config interface {
	InputFolderPath() string
	OutputFile() string
	ModuleName() string
	ModelsPkg() string
	ClearOutputDir() bool
}

var _ Config = &config{}

type config struct {
	inputFolderPath string
	outputFile      string
	moduleName      string
	modelsPkg       string
	clearOutputDir  bool
}

func NewConfig() *config {
	return &config{}
}

func (c *config) ClearOutputDir() bool {
	return c.clearOutputDir
}

func (c *config) InputFolderPath() string {
	return c.inputFolderPath
}

func (c *config) ModelsPkg() string {
	return c.modelsPkg
}

func (c *config) ModuleName() string {
	return c.moduleName
}

func (c *config) OutputFile() string {
	return c.outputFile
}

func (c *config) WithInputFolderPath(value string) error {
	var err error
	c.inputFolderPath, err = filepath.Abs(value)
	if err != nil {
		return err
	}

	// Check path exists and we have permission to read it
	if _, err := os.Stat(c.inputFolderPath); err != nil {
		return err
	}
	return nil
}

func (c *config) WithOutputFile(value string) error {
	var err error
	c.outputFile, err = filepath.Abs(value)
	if err != nil {
		return err
	}
	return nil
}

func (c *config) WithModuleName(value string) {
	c.moduleName = value
}

func (c *config) WithModelPkg(value string) {
	c.modelsPkg = value
}

func (c *config) WithClearOutputDir(value bool) {
	c.clearOutputDir = value
}
