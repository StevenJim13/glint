package basic

import (
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
)

func init() {
	// 匿名语言的规则集合，没有对应的语法解析功能，仅能分析文件内容、文件名、等信息，无法解析语法树。
	err := glint.InjectModels(utils.Any,
		&FileBasicModel,
	)

	errors.CheckErr(err)
	log.Infof("successfully injected %s models", utils.CCpp)
}
