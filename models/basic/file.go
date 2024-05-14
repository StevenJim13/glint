package basic

import (
	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/log"
)

var (
	charSetKey       = "charset"
	maxLinesKey      = "maxLines"
	maxLineLengthKey = "maxLineLength"
	newLineKey       = "newLine"
)

// FileBasic 与语言无关的规则
var FileBasic = glint.Model{
	Name: "FileBasic",
	Tags: []string{"basic"},
	Options: map[string]any{
		charSetKey:       "utf-8",
		maxLinesKey:      1200,
		maxLineLengthKey: 120,
		newLineKey:       "\\n",
	},
	ModelFunc: func(model *glint.Model, ctx glint.Context) {
		log.Infof("apply model %q to %s", model.Name, ctx.File())
		var ctt []byte
		// verify charset
		if value, ok := model.Options[charSetKey]; ok {
			charset := value.(string)
			ctt = ctx.Content()
			if err := verifyCharset(ctt, charset); err != nil {
				ctx.Defect(&model.Name, 0, 0, "the file charser expect %q", charset)
			}
		}

		var info glint.LinesInfo

		// verify MaxLines
		if value, ok := model.Options[maxLinesKey]; ok {
			lines := value.(int)
			if info == nil {
				info = ctx.LinesInfo()
			}
			if info.Lines() > lines {
				ctx.Defect(
					&model.Name,
					0, 0,
					"the lines count should be <= %d, but %d", info.Lines(),
				)
			}
		}

		// verify line max length
		if value, ok := model.Options[maxLineLengthKey]; ok {
			maxCharCount := value.(int)
			if info == nil {
				info = ctx.LinesInfo()
			}
			log.Infof("config limit line max length: %d", maxCharCount)
			info.Range(func(index int, line [2]int) bool {
				if line[0] > maxCharCount {
					log.Infof("add defect ... max file : %s", ctx.File())
					ctx.Defect(
						&model.Name,
						index, maxCharCount,
						"Expected %d characters or less per line, but found %d", maxCharCount, line[0],
					)
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
				info.Range(func(index int, line [2]int) bool {
					if line[1] != char {
						ctx.Defect(
							&model.Name,
							index,
							line[0],
							"Expected newline character is '\\n' not '\\r'",
						)
					}
					return true
				})
			}
		}
	},
}

// verifyCharset verify charset of file
func verifyCharset(content []byte, charset string) error {
	return nil
}
