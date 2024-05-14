/*
Copyright © 2024 clarkmonkey@163.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/stkali/glint/config"
	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "print default configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		err := Configure(configFile)
		errors.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// Configure ...
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

// configure ...
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

// generateDefaultConfig ...
func generateDefaultConfig() (*config.Config, error) {

	conf := &config.Config{
		Version:        config.Version,
		Concurrecy:     1024,
		LogLevel:       "error",
		LogFile:        "",
		WarningDisable: false,
		ResultFormat:   "cmd",
		ExcludeDirs:    []string{".*", "testdata"},
		ExcludeFiles:   []string{".*"},
	}
	modelSet := glint.ExportAllModels()
	conf.Languages = make([]*config.Language, 0, len(modelSet))
	for lang, modelList := range modelSet {

		modelCount := len(modelList)
		if modelCount == 0 {
			continue
		}

		exts, err := utils.Extends(lang)
		if err != nil {
			return nil, err
		}

		language := &config.Language{
			Name:    lang.String(),
			Extends: exts,
			Models:  make([]config.Model, 0, modelCount),
		}
		for _, model := range modelList {
			confMod := config.Model{
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
