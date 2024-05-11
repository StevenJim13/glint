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
	"encoding/csv"
	"strings"

	"io/fs"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stkali/glint/config"
	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/errors"
	"gopkg.in/yaml.v3"
)

const defaultConfigFile = "glint.yaml"

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var project string
		if len(args) == 0 {
			project = "."
		} else {
			project = args[0]
		}
		conf, err := getConfig(cmd.Flags())
		errors.CheckErr(err)
		err = glint.Lint(conf, project)
		errors.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringP("result-format", "", "cmd", "specify glint report format")
	lintCmd.Flags().StringP("result-file", "", "", "specify glint report file")
	lintCmd.Flags().IntP("concurrency", "c", 1024, "specify number of glint concurrency checker goroutine")
	lintCmd.Flags().BoolP("disable-warning", "", false, "specify number of glint concurrency checker goroutine")
	lintCmd.Flags().StringSliceP("exclude-tags", "", nil, "specify enable tag")
	lintCmd.Flags().StringSliceP("exclude-names", "", nil, "specify enable tag")
}

// getConfig ...
func getConfig(flags *pflag.FlagSet) (*config.Config, error) {

	if configFile == "" {
		configFile = defaultConfigFile
	}

	fd, err := os.OpenFile(configFile, os.O_RDONLY, os.ModePerm)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.Newf("not found config file: %q", configFile)
	}
	if err != nil {
		return nil, errors.Newf("failed to read config file: %q, err: %s", configFile, err)
	}

	dec := yaml.NewDecoder(fd)
	conf := config.NewConfig()
	if err = dec.Decode(conf); err != nil {
		return nil, errors.Newf("failed to unserialize config-file:%q, err: %s", configFile, err)
	}

	var nerr error
	flags.Visit(func(p *pflag.Flag) {
		switch p.Name {
		case "concurrency":
			conf.Concurrecy, nerr = strconv.Atoi(p.Value.String())
		case "result-format":
			conf.ResultFormat = p.Value.String()
		case "result-file":
			conf.ResultFile = p.Value.String()
		case "disable-warning":
			conf.WarningDisable, nerr = strconv.ParseBool(p.Value.String())
		case "exclude-tags":
			conf.ExcludeTags, nerr = parseStringSlice(p.Value.String())
		case "exclude-names":
			conf.ExcludeNames, nerr = parseStringSlice(p.Value.String())
		}
		err = errors.Join(err, nerr)
	})
	return conf, err
}

func parseStringSlice(s string) ([]string, error) {
	p := s[1 : len(s)-1]
	if p == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(p)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}
