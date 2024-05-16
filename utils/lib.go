package utils

import (
	"io"
	"io/fs"
	"time"

	"github.com/stkali/utility/errors"
)

func MatchNewChar(char string) int {
	switch char {
	case "\\r":
		return 0
	case "\\n":
		return 1
	case "\\r\\n":
		return 2
	default:
		return 100
	}
}

func Close(f any) {
	if closer, ok := f.(io.Closer); ok {
		var err error
		for i := 0; i < 3; i++ {

			if err = closer.Close(); err == nil {
				return
			} else if err != nil {
				if errInstance, ok := err.(*fs.PathError); ok && errInstance.Err.Error() == "file already closed" {
					return
				} else {
					time.Sleep(time.Millisecond * 200)
				}
			}
		}
		errors.CheckErr(err)
	}
}
