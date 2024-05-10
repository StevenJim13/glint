package config

import (
	"github.com/stkali/utility/errors"
)

const (
	Program = "glint"
	Version = "0.1.0"
)

func init() {
	errors.SetWarningPrefixf("%s warning", Program)
	errors.SetErrPrefixf("%s error", Program)
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

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		LogLevel: "error",
	}
}

// Validate ...
func Validate(conf *Config) error {
	if conf.Concurrecy < 1 {
		return errors.Newf("inavlid concurrency:%q must be > 0", conf.Concurrecy)
	}
	if err := IsCompatible(conf.Version); err != nil {
		return err
	}
	return nil
}

// IsCopatible ...
func IsCompatible(version string) error {
	marjor, err := getMarjor(version)
	if err != nil {
		return err
	}
	if marjor != Version[:len(marjor)] {
		return errors.Newf("incompatible version: %q, please upgrade to %s", marjor, version)
	}
	return nil
}

// getMarjor ...
func getMarjor(version string) (string, error) {

	if version == "" {
		return "nil", errors.New("empty version number")
	}

	for index, char := range version {
		if char > '9' || char < '0' {
			if index == 0 {
				return "", errors.Newf("invalid version: %q", version)
			}
		}
		return version[:index], nil
	}

	// only marjor number
	return version, nil
}
