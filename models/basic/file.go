package basic

import (
	"fmt"

	"github.com/stkali/glint/models"
	"github.com/stkali/glint/parser"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/log"
)

var (
	charSetKey       = "charset"
	maxLinesKey      = "maxLines"
	maxLineLengthKey = "maxLineLength"
	newLineKey       = "newLine"
)

// FileBasic ...
var FileBasic = models.Model{
	Name: "FileBasic",
	Tags: []string{"basic"},
	Options: map[string]any{
		charSetKey:       "utf-8",
		maxLinesKey:      1200,
		maxLineLengthKey: 120,
		newLineKey:       "\\n",
	},
	ModelFunc: func(model *models.Model, ctx parser.Context) {

		var ctt []byte
		// verify charset
		if value, ok := model.Options[charSetKey]; ok {
			charset := value.(string)
			ctt = ctx.Content()
			if err := verifyCharset(ctt, charset); err != nil {
				ctx.AddDefect(models.NewDefect(fmt.Sprintf("the file charset expect %q", charset)))
			}
		}

		var info parser.LinesInfo

		// verify MaxLine
		if value, ok := model.Options[maxLinesKey]; ok {
			lines := value.(int)
			if info == nil {
				info = ctx.LinesInfo()
			}
			if info.Lines() > lines {
				ctx.AddDefect(
					models.NewDefect(fmt.Sprintf("the lines count should be <= %d, but %d", lines, info.Lines())),
				)
			}
		}

		// verify line max length
		if value, ok := model.Options[maxLineLengthKey]; ok {
			length := value.(int)
			if info == nil {
				info = ctx.LinesInfo()
			}
			info.Range(func(line [2]int) bool {
				if line[0] > int(length) {
					ctx.AddDefect(models.NewDefect(" too big"))
				}
				return true
			})
		}

		// verify new line character
		if value, ok := model.Options[newLineKey]; ok {
			newLineChar := value.(string)
			char := utils.MatchNewChar(newLineChar)
			if info == nil {
				info = ctx.LinesInfo()
				info.Range(func(line [2]int) bool {
					if line[1] != char {
						ctx.AddDefect(models.NewDefect(""))
					}
					return true
				})
			}
		}
	},
}

// verifyCharset verify charset of file
func verifyCharset(content []byte, charset string) error {
	log.Infof("verify charset %q", charset)
	return nil
}
