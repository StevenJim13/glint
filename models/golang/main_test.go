package golang

import (
	"testing"

	"github.com/stkali/utility/log"
	"github.com/stkali/utility/tool"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.INFO)
	tool.Exit(m.Run())
}
