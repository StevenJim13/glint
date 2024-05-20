package basic

import (
	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

const (
	charsetKey       = "charset"
	maxLinesKey      = "maxLines"
	maxLineLengthKey = "maxLineLength"
	newLineKey       = "newLine"
)

// FileBasic 与语言无关的规则
var FileBasicModel = glint.Model{
	Name: "FileBasic",
	Tags: []string{"basic"},
	Options: map[string]any{
		charsetKey:       "utf-8",
		maxLinesKey:      1200,
		maxLineLengthKey: 120,
		newLineKey:       "\\n",
	},
	GenerateModelFunc: func(model *glint.Model) (glint.CheckFuncType, error) {

		var charset *string
		v, ok := model.Options[charsetKey]
		if ok {
			set, ok := v.(string)
			if !ok {
				return nil, errors.Newf("%s want string but get %q", charsetKey, v)
			} else {
				charset = &set
			}
		}

		var lines *int
		v, ok = model.Options[maxLinesKey]
		if ok {
			if l, ok := v.(int); !ok {
				return nil, errors.Newf("%s want int but get %q", maxLinesKey, v)
			} else {
				lines = &l
			}
		}

		var length *int
		v, ok = model.Options[maxLineLengthKey]
		if ok {
			if l, ok := v.(int); !ok {
				return nil, errors.Newf("%s want int but get %q", maxLineLengthKey, v)
			} else {
				length = &l
			}
		}

		var newline *int
		var newlineChar string
		v, ok = model.Options[newLineKey]
		if ok {
			if newlineChar, ok = v.(string); !ok {
				return nil, errors.Newf("%s want int but get %q", newLineKey, v)
			} else {
				t := utils.NewLineType(newlineChar)
				newline = &t
			}
		}

		return func(ctx glint.Context) error {
			log.Infof("apply %s model", model.Name)
			var content []byte
			if charset != nil {
				content = ctx.Content()
				if actual, ok := utils.VerifyCharset(content, *charset); !ok {
					glint.AddDefect(ctx, model, 0, 0,
						"The expected charset is %q not %q", *charset, actual)
				}
			}

			var info glint.LinesInfo
			if lines != nil {
				if info == nil {
					info = ctx.Lines()
				}
				if info.Lines() > *lines {
					glint.AddDefect(ctx, model, 0, 0,
						"Expected %d characters or less per file, but found %d", lines, info.Lines())
				}
			}

			if length != nil && newline != nil {
				if info == nil {
					info = ctx.Lines()
				}
				info.Range(func(index int, line [2]int) bool {
					if line[0] > *length {
						glint.AddDefect(ctx, model, index, line[0],
							"Expected %d characters or less per line, but found %d", *length, line[0])
					}
					if line[1] != *newline {
						if line[1] == -1 {
							glint.AddDefect(ctx, model, index, line[0],
								"Expected to find a blank line at the end of the file, but didn't")
						} else {
							glint.AddDefect(ctx, model, index, line[0],
								"Expected newline character is %q not %q", newlineChar, utils.NewLineChar(line[1]))
						}
					}
					return true
				})
			}

			return nil
		}, nil
	},
}
