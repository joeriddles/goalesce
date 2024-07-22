package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/generate"
	"github.com/joeriddles/gorm-oapi-codegen/pkg/parse"
)

var (
	flagPrintUsage bool

	flagOutputFile string
	flagModuleName string
	flagModelsPkg  string
)

func main() {
	flag.BoolVar(&flagPrintUsage, "help", false, "Show this help and exit.")
	flag.BoolVar(&flagPrintUsage, "h", false, "Same as -help.")

	flag.StringVar(&flagOutputFile, "o", "./generated", "Where to output generated code, ./generated/ is default.")
	flag.StringVar(&flagModuleName, "module", "", "The name of the module the generated code will be part of")
	flag.StringVar(&flagModelsPkg, "pkg", "", "The name of the package that the GORM models are part of")

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
	if err := run(folderPath, flagOutputFile, flagModuleName, flagModelsPkg); err != nil {
		errExit(err.Error())
	}
}

func run(folderPath string, outputPath, moduleName, modelsPkgName string) error {
	folderPath, err := filepath.Abs(folderPath)
	if err != nil {
		return err
	}

	outputPath, err = filepath.Abs(outputPath)
	if err != nil {
		return nil
	}

	// Check path exists and we have permission to read it
	if _, err := os.Stat(folderPath); err != nil {
		return err
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	parser := parse.NewParser()
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

		generator := generate.NewGenerator(outputPath, moduleName, modelsPkgName)
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
