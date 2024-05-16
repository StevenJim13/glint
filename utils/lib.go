package utils

import (
	"io"
	"io/fs"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/stkali/utility/errors"
	"golang.org/x/net/html/charset"
)

func NewLineType(char string) int {
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

func NewLineChar(t int) string {
	switch t {
	case 0:
		return "\\r"
	case 1:
		return "\\n"
	case 2:
		return "\\r\\n"
	default:
		return ""
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

func IsPureAscii(s string) bool {
	for _, char := range s {
		if char > 256 {
			return false
		}
	}
	return true
}

// verifyCharset verify charset of file
// "utf-8",
// "ibm866",
// "iso-8859-2",
// "iso-8859-3",
// "iso-8859-4",
// "iso-8859-5",
// "iso-8859-6",
// "iso-8859-7",
// "iso-8859-8",
// "iso-8859-8-i",
// "iso-8859-10",
// "iso-8859-13",
// "iso-8859-14",
// "iso-8859-15",
// "iso-8859-16",
// "koi8-r",
// "koi8-u",
// "macintosh",
// "windows-874",
// "windows-1250",
// "windows-1251",
// "windows-1252",
// "windows-1253",
// "windows-1254",
// "windows-1255",
// "windows-1256",
// "windows-1257",
// "windows-1258",
// "x-mac-cyrillic",
// "gbk",
// "gb18030",
// "big5",
// "euc-jp",
// "iso-2022-jp",
// "shift_jis",
// "euc-kr",
// "replacement",
// "utf-16be",
// "utf-16le",
// "x-user-defined",
func VerifyCharset(content []byte, encoding string) (string, bool) {

	if content == nil {
		return encoding, true
	}
	lowEncoding := strings.ToLower(encoding)
	var DetCoding string
	switch lowEncoding {
	case "utf-8", "utf8":
		if utf8.Valid(content) {
			return encoding, true
		}
	case "gbk":
		if isGBK(content) {
			return encoding, true
		}
		_, DetCoding, _ = charset.DetermineEncoding(content, "text")
		return DetCoding, false
	default:
		_, DetCoding, _ = charset.DetermineEncoding(content, "text")
		if strings.ReplaceAll(lowEncoding, "-", "") == strings.ReplaceAll(DetCoding, "-", "") {
			return encoding, true
		}
		return DetCoding, false
	}
	if DetCoding != "" {
		_, DetCoding, _ = charset.DetermineEncoding(content, "text")
	}
	return DetCoding, false
}

func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}
