package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/config"
	"github.com/joeriddles/gorm-oapi-codegen/pkg/generate"
	"github.com/joeriddles/gorm-oapi-codegen/pkg/parse"
)

var (
	flagPrintUsage bool

	flagOutputFile        string
	flagModuleName        string
	flagModelsPkg         string
	flagClearOutputDir    bool
	flagAllowCustomModels bool
)

func main() {
	flag.BoolVar(&flagPrintUsage, "help", false, "Show this help and exit.")
	flag.BoolVar(&flagPrintUsage, "h", false, "Same as -help.")

	flag.StringVar(&flagOutputFile, "o", "./generated", "Where to output generated code, ./generated/ is default.")
	flag.StringVar(&flagModuleName, "module", "", "The name of the module the generated code will be part of")
	flag.StringVar(&flagModelsPkg, "pkg", "", "The name of the package that the GORM models are part of")
	flag.BoolVar(&flagClearOutputDir, "clear", false, "If true, clears the contents of the output directory before generating new files")
	flag.BoolVar(&flagAllowCustomModels, "custom", false, "If true, parses classes that do not inherit from gorm.Model")

	flag.Parse()

	if flagPrintUsage {
		flag.Usage()
		os.Exit(0)
	}

	if flagModuleName == "" {
		errExit("Please specify a module name with -module\n")
	}

	if flagModelsPkg == "" {
		errExit("Please specify a package name for the GORM models with -pkg\n")
	}

	if flag.NArg() < 1 {
		errExit("Please specify a path to a folder of GORM models\n")
	} else if flag.NArg() > 1 {
		errExit("Only one folder path is accepted and it must be the last CLI argument\n")
	}

	folderPath := flag.Arg(0)

	cfg := config.NewConfig()

	err := cfg.WithInputFolderPath(folderPath)
	if err != nil {
		errExit("Invalid input folder path: %v\n", folderPath)
	}

	err = cfg.WithOutputFile(flagOutputFile)
	if err != nil {
		errExit("Invalid output filepath: %v\n", flagOutputFile)
	}

	cfg.WithModuleName(flagModuleName)
	cfg.WithModelPkg(flagModelsPkg)
	cfg.WithClearOutputDir(flagClearOutputDir)

	if err := run(cfg); err != nil {
		errExit(err.Error())
	}
}

func run(cfg config.Config) error {
	folderPath := cfg.InputFolderPath()
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	parser := parse.NewParser(log.Default(), cfg.AllowCustomModels())
	for _, entry := range entries {
		filename := entry.Name()
		if !strings.HasSuffix(filename, ".go") {
			continue
		}

		entryFilepath := filepath.Join(folderPath, filename)
		metadatas, err := parser.Parse(entryFilepath)
		if err != nil {
			return err
		}

		generator, err := generate.NewGenerator(
			cfg.OutputFile(),
			cfg.ModuleName(),
			cfg.ModelsPkg(),
			cfg.ClearOutputDir(),
		)
		if err != nil {
			return err
		}
		if err := generator.Generate(metadatas); err != nil {
			return err
		}
	}

	return nil
}

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
