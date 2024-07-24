package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	InputFolderPath   string `yaml:"input_folder_path"`
	OutputFile        string `yaml:"output_file_path"`
	ModuleName        string `yaml:"module_name"`
	ModelsPkg         string `yaml:"models_package"`
	ClearOutputDir    bool   `yaml:"clear_output_dir"`
	AllowCustomModels bool   `yaml:"allow_custom_models"`
	PruneYaml         bool   `yaml:"prune_yaml"`
}

func (o *Config) Validate() error {
	var errs []error

	if o.ModuleName == "" {
		errs = append(errs, errors.New("module_name must be specified"))
	}

	if o.ModelsPkg == "" {
		errs = append(errs, errors.New("model_package must be specified"))
	}

	if o.OutputFile == "" {
		errs = append(errs, errors.New("output_file_path must be specified"))
	}

	err := errors.Join(errs...)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %w", err)
	}

	o.InputFolderPath, err = filepath.Abs(o.InputFolderPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(o.InputFolderPath); err != nil {
		return err
	}

	o.OutputFile, err = filepath.Abs(o.OutputFile)
	if err != nil {
		return err
	}

	return nil
}
