package config

import (
	"fmt"

	"github.com/stkali/utility/errors"
)

const (
	Program = "glint"
	Version = "0.1.0"
)

type Config struct {
	Version        string      `yaml:"version"`
	Concurrecy     int         `yaml:"concurrecy"`
	WarningDisable bool        `yaml:"warningDisable"`
	OutputFile     string      `yaml:"resultFile"`
	OutputFormat   string      `yaml:"resultFormat"`
	ExcludeTags    []string    `yaml:"excludeTags`
	ExcludeNames   []string    `yaml:"excludeNames`
	ExcludeDirs    []string    `yaml:"excludeDirs"`
	ExcludeFiles   []string    `yaml:"excludeFiles"`
	Languages      []*Language `yaml:"languages"`
}

// Validate validates version compatibility and legality of settings.
func Validate(conf *Config) error {

	if conf.Concurrecy < 1 {
		return errors.Newf("Inavlid concurrency:%d must be > 0", conf.Concurrecy)
	}

	if err := IsCompatible(conf.Version); err != nil {
		return err
	}
	for _, lang := range conf.Languages {
		if lang.Name == "" {
			return errors.New("Language is empty")
		}
		for _, model := range lang.Models {
			if model.Name == "" {
				return errors.New("Model name is empty")
			}
		}
	}
	return nil
}

// IsCopatible verifies version number compatibility.
func IsCompatible(version string) error {
	marjor, err := getMarjor(version)
	if err != nil {
		return err
	}
	if *marjor != Version[:len(*marjor)] {
		return errors.Newf("Incompatible version: %q, please upgrade to %s", marjor, version)
	}
	return nil
}

// String implements 'fmt.Stringer' interface.
func (c *Config) String() string {
	return fmt.Sprintf("<Config version: %s>", c.Version)
}

// getMarjor returns the master version number if version is not empty else return error
func getMarjor(version string) (*string, error) {

	if version == "" {
		return nil, errors.New("Version number is empty")
	}

	for index, char := range version {
		if char > '9' || char < '0' {
			if index == 0 {
				return nil, errors.Newf("Invalid version number: %q", version)
			}
		}
		marjor := version[:index]
		return &marjor, nil
	}

	// only marjor number
	return &version, nil
}
