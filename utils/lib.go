package utils

import (
	"io"
	"io/fs"
	"reflect"
	"time"

	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
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
			err = closer.Close()
			err = closer.Close()
			if err != nil {
				if errInstance, ok := err.(*fs.PathError); ok {
					log.Infof("err: %+v", errInstance)
				}
			}
			log.Infof("err: %s, type: %s", err, reflect.TypeOf(err))
			if err == nil {
				break
			} else {
				time.Sleep(time.Millisecond * 200)
			}
		}
		errors.CheckErr(err)
	}
}
