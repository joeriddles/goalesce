package generator

import (
	"github.com/joeriddles/goalesce/pkg/v2/config"
	"github.com/joeriddles/goalesce/pkg/v2/logger"
)

// Wrapper class to encapsulate all services needed for a Generator.
type GeneratorServices struct {
	Config        *config.Config
	LoggerFactory logger.LoggerFactory
}

func NewGeneratorServices(config *config.Config) *GeneratorServices {
	return &GeneratorServices{
		Config:        config,
		LoggerFactory: logger.NewLoggerFactory(),
	}
}
