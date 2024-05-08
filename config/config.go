package config

import (
	"fmt"
	"io"
	"os"

	"github.com/stkali/glint/models"
	"github.com/stkali/glint/util"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/tool"
	"gopkg.in/yaml.v3"
)

const (
	Program = "glint"
	Version = "0.1.0"
)

func init() {
	errors.SetWarningPrefix(fmt.Sprintf("%s warning:", Program))
	tool.SetErrorPrefix(fmt.Sprintf("%s error:", Program))
}

func IsCompatible(newVersion, oldVersion string) bool {
	return true
}

type Config struct {
	Version        string     `mapstructure:"version"`
	Concurrecy     int        `mapstructure:"concurrency"`
	LogLevel       string     `mapstructure:"logLevel"`
	LogFile        string     `mapstructure:"logFile"`
	WarningDisable bool       `mapstructure:"warningDisable"`
	ResultFile     string     `mapstructure:"resultFile"`
	ResultFormat   string     `mapstructure:"resultFormat"`
	ExcludeTags    []string   `mapstructure:"excludeTags`
	ExcludeNames   []string   `mapstructure:"excludeNames`
	Languages      []Language `mapstructure:"languages"`
}

func (c *Config) Validate() error {
	return nil
}

func NewConfig() *Config {
	return &Config{Version: Version, LogLevel: "error"}
}

func Configure(configPath string) error {
	var writer io.Writer
	if configPath == "" {
		writer = os.Stdout
	} else {
		f, err := os.OpenFile(configPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return errors.Newf("failed to create glint config file at: %q, err: %s", configPath, err)
		}
		defer f.Close()
		writer = f
	}
	return configure(writer)
}

func configure(writer io.Writer) error {

	conf, err := generateDefaultConfig()
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(writer)
	enc.SetIndent(2)
	defer enc.Close()
	err = enc.Encode(conf)
	if err != nil {
		return errors.Newf("failed to serialized config to yaml, err: %s", err)
	}
	return nil
}

func generateDefaultConfig() (*Config, error) {

	conf := &Config{
		Version:        Version,
		Concurrecy:     1024,
		LogLevel:       "info",
		LogFile:        "",
		WarningDisable: false,
		ResultFormat:   "cmd",
	}
	modelSet := models.ExportAllModels()
	conf.Languages = make([]Language, 0, len(modelSet))
	for lang, modelList := range modelSet {

		modelCount := len(modelList)
		if modelCount == 0 {
			continue
		}

		exts, err := util.Extends(lang)
		if err != nil {
			return nil, err
		}

		language := Language{
			Name:    lang.String(),
			Extends: exts,
			Models:  make([]Model, 0, modelCount),
		}
		for _, model := range modelList {
			confMod := Model{
				Name:    model.Name,
				Tags:    model.Tags,
				Options: model.Options,
			}
			language.Models = append(language.Models, confMod)
		}
		conf.Languages = append(conf.Languages, language)
	}
	return conf, nil
}
