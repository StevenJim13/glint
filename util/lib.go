package util

import (
	"fmt"
	"github.com/stkali/utility/tool"
	"github.com/stkali/utility/errors"
)

func SetWarningPrefixf(s string, args ... any) {
	errors.SetWarningPrefix(fmt.Sprintf(s, args...))
}

func SetErrorPrefixf(s string, args ... any) {
	tool.SetErrorPrefix(fmt.Sprintf(s, args...))
}

func IsCompatible(newVersion, oldVersion string) bool {
	return true
}