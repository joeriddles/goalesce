package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeriddles/goalesce/pkg/v2/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"gopkg.in/yaml.v2"
)

type Config struct {
	InputFolderPath string `yaml:"input_folder_path"`
	// Where to output generated code, ./generated/ is default
	OutputFile string `yaml:"output_file_path"`
	// The name of the module the generated code will be part of
	ModuleName string `yaml:"module_name"`
	// The name of the package that the GORM models are part of
	ModelsPkg string `yaml:"models_package"`
	// The name of the package that the GORM-generated Query is in
	QueryPkg string `yaml:"query_package"`
	// If true, clears the contents of the output directory before generating new files
	ClearOutputDir bool `yaml:"clear_output_dir"`
	// If true, parses classes that do not inherit from gorm.Model
	AllowCustomModels bool `yaml:"allow_custom_models"`
	// If true, deletes all model specific YAML files after combining them into a single YAML file
	PruneYaml bool `yaml:"prune_yaml"`
	// If true, the generated OpenAPI YAML file uses this as its base
	OpenApiFile string `yaml:"openapi_file"`
	// Excludes these GORM models from the generated OpenAPI routes
	ExcludeModels []string `yaml:"exclude_models,omitempty"`
	// If true, generates a sample main.go file for running the server
	GenerateMain bool `yaml:"generate_main"`
	// Override built-in templates from user-provided files
	UserTemplates map[string]string `yaml:"user_templates,omitempty"`

	// Generated Repository configuration
	RepositoryConfiguration *RepositoryConfiguration `yaml:"repository,omitempty"`

	// oapi-codegen server configuration
	ServerCodegen *OApiGenConfiguration `yaml:"server_codegen,omitempty"`
	// oapi-codegen types configuration
	TypesCodegen *OApiGenConfiguration `yaml:"types_codegen,omitempty"`
}

type OApiGenConfiguration struct {
	codegen.Configuration `yaml:",inline"`

	// OutputFile is the filepath to output
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

	if cfg.RepositoryConfiguration != nil {
		if isRelativeFilepath(cfg.RepositoryConfiguration.OutputFile) {
			cfg.RepositoryConfiguration.OutputFile = filepath.Join(configDir, cfg.RepositoryConfiguration.OutputFile)
		}
		if cfg.RepositoryConfiguration.Template != nil && isRelativeFilepath(*cfg.RepositoryConfiguration.Template) {
			templateFp := filepath.Join(configDir, *cfg.RepositoryConfiguration.Template)
			cfg.RepositoryConfiguration.Template = &templateFp
		}
	}

	if cfg.UserTemplates != nil {
		for key, path := range cfg.UserTemplates {
			cfg.UserTemplates[key] = filepath.Join(configDir, path)
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	var errs []error

	if c.ModuleName == "" {
		errs = append(errs, errors.New("module_name must be specified"))
	}

	if c.ModelsPkg == "" {
		errs = append(errs, errors.New("model_package must be specified"))
	}

	if c.QueryPkg == "" {
		errs = append(errs, errors.New("query_package must be specified"))
	}

	if c.OutputFile == "" {
		errs = append(errs, errors.New("output_file_path must be specified"))
	}

	err := errors.Join(errs...)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %w", err)
	}

	c.InputFolderPath, err = filepath.Abs(c.InputFolderPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(c.InputFolderPath); err != nil {
		return err
	}

	c.OutputFile, err = filepath.Abs(c.OutputFile)
	if err != nil {
		return err
	}

	if c.OpenApiFile != "" {
		c.OpenApiFile, err = filepath.Abs(c.OpenApiFile)
		if err != nil {
			return err
		}
	}

	if c.TypesCodegen == nil {
		outputFile := filepath.Join(c.OutputFile, "api", "types.gen.go")
		c.TypesCodegen = &OApiGenConfiguration{
			OutputFile: outputFile,
			Configuration: codegen.Configuration{
				PackageName: "api",
				Generate:    codegen.GenerateOptions{Models: true},
			},
		}
	}

	if c.ServerCodegen == nil {
		outputFile := filepath.Join(c.OutputFile, "api", "server_interface.gen.go")
		c.ServerCodegen = &OApiGenConfiguration{
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

	if c.RepositoryConfiguration == nil {
		c.RepositoryConfiguration = &RepositoryConfiguration{
			OutputFile: filepath.Join(c.OutputFile, "repository"),
		}
	}
	if err := c.RepositoryConfiguration.Validate(); err != nil {
		return err
	}

	return nil
}

// Get the path to the config's go.mod file.
func (c *Config) GetModulePath() (string, error) {
	modulePath, err := utils.FindGoMod(c.OutputFile, c.ModuleName)
	if err != nil {
		return "", err
	}
	return modulePath, nil
}

// Get the path to the folder containing the config's go.mod file.
func (c *Config) GetModuleFolderPath() (string, error) {
	modulePath, err := c.GetModulePath()
	if err != nil {
		return "", err
	}
	return filepath.Dir(modulePath), nil
}

// Get the relative path to the output folder from the go.mod folder path.
func (c *Config) GetRelativeOutputPath() (string, error) {
	moduleRootPath, err := c.GetModuleFolderPath()
	if err != nil {
		return "", err
	}

	outputPath, err := filepath.Rel(moduleRootPath, c.OutputFile)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

// Get the package name for the generated oapi-codegen types code.
func (c *Config) GetTypesPackage() (string, error) {
	moduleRootPath, err := c.GetModuleFolderPath()
	if err != nil {
		return "", err
	}

	relPath, err := c.GetRelativeOutputPath()
	if err != nil {
		return "", err
	}

	var typesPackage string = filepath.Join(c.ModuleName, relPath, "api")
	defaultOutputFile := filepath.Join(c.OutputFile, "api", "types.gen.go")
	if c.TypesCodegen.OutputFile != defaultOutputFile {
		relPkg, err := filepath.Rel(moduleRootPath, c.TypesCodegen.OutputFile)
		if err != nil {
			return "", err
		}
		pkg := filepath.Join(c.ModuleName, relPkg)
		pkg = filepath.Dir(pkg) // remove filename
		typesPackage = pkg
	}

	return typesPackage, nil
}

// Get the package name for the generated repository code.
func (c *Config) GetRepositoryPackage() (string, error) {
	moduleRootPath, err := c.GetModuleFolderPath()
	if err != nil {
		return "", err
	}

	relPath, err := c.GetRelativeOutputPath()
	if err != nil {
		return "", err
	}

	var repositoryPackage string = filepath.Join(c.ModuleName, relPath, "repository")
	defaultRepositoryOutputFile := filepath.Join(c.OutputFile, "repository")
	if c.RepositoryConfiguration.OutputFile != defaultRepositoryOutputFile {
		relPkg, err := filepath.Rel(moduleRootPath, c.RepositoryConfiguration.OutputFile)
		if err != nil {
			return "", err
		}
		repositoryPackage = filepath.Join(c.ModuleName, relPkg)
	}
	return repositoryPackage, nil
}

// Get the package name for the GORM model code.
func (c *Config) GetModelPackage() (string, error) {
	return c.ModelsPkg, nil
}

// Get the package name for the GORM query code.
func (c *Config) GetQueryPackage() (string, error) {
	return c.QueryPkg, nil
}

type RepositoryConfiguration struct {
	// OutputFile is the folder to output to
	OutputFile string `yaml:"output,omitempty"`
	// Custom generated repository template
	Template *string `yaml:"template,omitempty"`
}

func (c *RepositoryConfiguration) Validate() error {
	var errs []error
	if c.OutputFile == "" {
		errs = append(errs, errors.New("output must be specified"))
	}
	err := errors.Join(errs...)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %w", err)
	}
	return nil
}

func isRelativeFilepath(fp string) bool {
	return !strings.HasPrefix("/", fp)
}
