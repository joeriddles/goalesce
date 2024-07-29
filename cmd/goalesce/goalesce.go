package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/joeriddles/goalesce/pkg/generate"
	"github.com/joeriddles/goalesce/pkg/parse"
	"golang.org/x/tools/go/packages"
)

var (
	flagPrintUsage bool

	flagConfigFile        string
	flagOutputFile        string
	flagModuleName        string
	flagModelsPkg         string
	flagClearOutputDir    bool
	flagAllowCustomModels bool
	flagPruneYaml         bool
)

const LoadAll = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax

func main() {
	flag.BoolVar(&flagPrintUsage, "help", false, "Show this help and exit.")
	flag.BoolVar(&flagPrintUsage, "h", false, "Same as -help.")

	flag.StringVar(&flagConfigFile, "config", "", "A YAML config file that controls goalesce behavior.")
	flag.StringVar(&flagOutputFile, "o", "./generated", "Where to output generated code, ./generated/ is default.")
	flag.StringVar(&flagModuleName, "module", "", "The name of the module the generated code will be part of")
	flag.StringVar(&flagModelsPkg, "pkg", "", "The name of the package that the GORM models are part of")
	flag.BoolVar(&flagClearOutputDir, "clear", false, "If true, clears the contents of the output directory before generating new files")
	flag.BoolVar(&flagAllowCustomModels, "custom", false, "If true, parses classes that do not inherit from gorm.Model")
	flag.BoolVar(&flagPruneYaml, "prune", false, "If true, deletes all model specific YAML files after combining them into a single YAML file")

	flag.Parse()

	if flagPrintUsage {
		flag.Usage()
		os.Exit(0)
	}

	var cfg *config.Config
	var err error
	if flagConfigFile != "" {
		cfg, err = config.FromYamlFile(flagConfigFile)
		if err != nil {
			errExit("%v", err)
		}
	}

	if cfg == nil {
		if flag.NArg() < 1 {
			errExit("Please specify a path to a folder of GORM models\n")
		} else if flag.NArg() > 1 {
			errExit("Only one folder path is accepted and it must be the last CLI argument\n")
		}

		cfg = &config.Config{
			InputFolderPath:   flag.Args()[0],
			OutputFile:        flagOutputFile,
			ModuleName:        flagModuleName,
			ModelsPkg:         flagModelsPkg,
			ClearOutputDir:    flagClearOutputDir,
			AllowCustomModels: flagAllowCustomModels,
			PruneYaml:         flagPruneYaml,
		}
	}

	if err := cfg.Validate(); err != nil {
		errExit("configuration error: %v\n", err)
	}

	if err := Run(cfg); err != nil {
		errExit(err.Error())
	}
}

func Run(cfg *config.Config) error {
	logger := log.Default()
	parser := parse.NewParser(logger, cfg)

	metadatas, err := parser.Parse(cfg.InputFolderPath)
	if err != nil {
		return err
	}

	generator, err := generate.NewGenerator(logger, cfg)
	if err != nil {
		return err
	}
	if err := generator.Generate(metadatas); err != nil {
		return err
	}

	return nil
}

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
