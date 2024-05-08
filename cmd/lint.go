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
	"os"

	"github.com/spf13/cobra"
	"github.com/stkali/glint/config"
	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
	"github.com/stkali/utility/tool"
	"gopkg.in/yaml.v3"
)

var c = config.NewConfig()

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
		conf, err := getConfig()
		tool.CheckError("failed get config, err: %s", err)
		glint.Lint(conf, project)
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
	// result form
	lintCmd.Flags().StringVarP(&c.ResultFile, "result-format", "", "cmd", "specify glint report format")
	// result file
	lintCmd.Flags().StringVarP(&c.ResultFormat, "result-file", "", "", "specify glint report file")
	// cocurrency
	lintCmd.Flags().IntVarP(&c.Concurrecy, "concurrency", "", 1024, "specify number of glint concurrency checker goroutine")
	lintCmd.Flags().StringArrayVarP(&c.ExcludeTags, "exclude-tags", "", nil, "specify enable tag")
	lintCmd.Flags().StringArrayVarP(&c.ExcludeNames, "exclude-names", "", nil, "specify enable tag")
}

func getConfig() (*config.Config, error) {
	
	fd, err := os.OpenFile(configFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, errors.Newf("failed to read config file: %s, err: %s", configFile, err)
	}
	
	dec := yaml.NewDecoder(fd)
	conf := config.NewConfig()
	if err = dec.Decode(conf); err != nil {
		return nil, errors.Newf("failed to unserialize config-file:%q, err: %s", configFile, err)
	}
	
	rootCmd.Commands()
	log.Info(lintCmd.Flags().GetString("result-format"))
	return conf, nil
}
