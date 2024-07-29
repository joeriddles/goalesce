package pkg

import (
	"log"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/joeriddles/goalesce/pkg/generate"
	"github.com/joeriddles/goalesce/pkg/parse"
)

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
