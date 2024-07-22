package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/generate"
	"github.com/joeriddles/gorm-oapi-codegen/pkg/parse"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		errExit("Please specify a path to a folder of GORM models\n")
	} else if flag.NArg() > 1 {
		errExit("Only one folder path is accepted and it must be the last CLI argument\n")
	}

	folderPath := flag.Arg(0)
	if err := run(folderPath); err != nil {
		errExit(err.Error())
	}
}

func run(folderPath string) error {
	folderPath, err := filepath.Abs(folderPath)
	if err != nil {
		return err
	}

	// Check path exists and we have permission to read it
	if _, err := os.Stat(folderPath); err != nil {
		return err
	}

	entries, err := os.ReadDir(folderPath)
	for _, entry := range entries {
		filename := entry.Name()
		entryFilepath := filepath.Join(folderPath, filename)
		metadatas, err := parse.Parse(entryFilepath)
		if err != nil {
			return err
		}
		if err := generate.Generate(metadatas); err != nil {
			return err
		}
	}
	return err
}

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
