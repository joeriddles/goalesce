package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"gopkg.in/yaml.v2"
)

type Config struct {
	InputFolderPath   string                `yaml:"input_folder_path"`
	OutputFile        string                `yaml:"output_file_path"`
	ModuleName        string                `yaml:"module_name"`
	ModelsPkg         string                `yaml:"models_package"`
	ClearOutputDir    bool                  `yaml:"clear_output_dir"`
	AllowCustomModels bool                  `yaml:"allow_custom_models"`
	PruneYaml         bool                  `yaml:"prune_yaml"`
	OpenApiFile       string                `yaml:"openapi_file"`
	ServerCodegen     *OApiGenConfiguration `yaml:"server_codegen,omitempty"`
	TypesCodegen      *OApiGenConfiguration `yaml:"types_codegen,omitempty"`
	ExcludeModels     []string              `yaml:"exclude_models,omitempty"`
	GenerateMain      bool                  `yaml:"generate_main"`
}

type OApiGenConfiguration struct {
	codegen.Configuration `yaml:",inline"`

	// OutputFile is the filepath to output.
	OutputFile string `yaml:"output,omitempty"`
}

func FromYamlFile(fp string) (*Config, error) {
	absoluteConfigFile, err := filepath.Abs(fp)
	if err != nil {
		return nil, fmt.Errorf("config file not found '%s': %v", fp, err)
	}
	configFile, err := os.ReadFile(absoluteConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file '%s': %v", fp, err)
	}
	cfg := &Config{}
	if err = yaml.UnmarshalStrict(configFile, cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	// Make any relative paths relative to the YAML config filepath
	configDir := filepath.Dir(absoluteConfigFile)
	if isRelativeFilepath(cfg.InputFolderPath) {
		cfg.InputFolderPath = filepath.Join(configDir, cfg.InputFolderPath)
	}
	if isRelativeFilepath(cfg.OutputFile) {
		cfg.OutputFile = filepath.Join(configDir, cfg.OutputFile)
	}
	if isRelativeFilepath(cfg.OpenApiFile) {
		cfg.OpenApiFile = filepath.Join(configDir, cfg.OpenApiFile)
	}
	if cfg.ServerCodegen != nil && isRelativeFilepath(cfg.ServerCodegen.OutputFile) {
		cfg.ServerCodegen.OutputFile = filepath.Join(configDir, cfg.ServerCodegen.OutputFile)
	}
	if cfg.TypesCodegen != nil && isRelativeFilepath(cfg.TypesCodegen.OutputFile) {
		cfg.TypesCodegen.OutputFile = filepath.Join(configDir, cfg.TypesCodegen.OutputFile)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
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

	if o.OpenApiFile != "" {
		o.OpenApiFile, err = filepath.Abs(o.OpenApiFile)
		if err != nil {
			return err
		}
	}

	if o.TypesCodegen == nil {
		outputFile := filepath.Join(o.OutputFile, "api", "types.gen.go")
		o.TypesCodegen = &OApiGenConfiguration{
			OutputFile: outputFile,
			Configuration: codegen.Configuration{
				PackageName: "api",
				Generate:    codegen.GenerateOptions{Models: true},
			},
		}
	}

	if o.ServerCodegen == nil {
		outputFile := filepath.Join(o.OutputFile, "api", "server_interface.gen.go")
		o.ServerCodegen = &OApiGenConfiguration{
			OutputFile: outputFile,
			Configuration: codegen.Configuration{
				PackageName: "api",
				Generate: codegen.GenerateOptions{
					StdHTTPServer: true,
					Strict:        true,
					EmbeddedSpec:  true,
				},
			},
		}
	}

	return nil
}

func isRelativeFilepath(fp string) bool {
	return !strings.HasPrefix("/", fp)
}
