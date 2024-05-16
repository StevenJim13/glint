package basic

import (
	"strings"
	"unicode/utf8"

	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
	"golang.org/x/net/html/charset"
)

var (
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
	ModelFunc: func(model *glint.Model, ctx glint.Context) {

		var ctt []byte
		// verify charset
		if value, ok := model.Options[charsetKey]; ok {
			charset := value.(string)
			ctt = ctx.Content()
			if actual, ok := verifyCharset(ctt, charset); !ok {
				ctx.Defect(model, 0, 0,
					"The expected charset is %q not %q", charset, actual)
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
				ctx.Defect(model, 0, 0,
					"Expected %d characters or less per file, but found %d", lines, info.Lines(),
				)
			}
		}

		// verify line max length
		if value, ok := model.Options[maxLineLengthKey]; ok {
			maxCharCount := value.(int)
			if info == nil {
				info = ctx.LinesInfo()
			}
			info.Range(func(index int, line [2]int) bool {
				if line[0] > maxCharCount {
					ctx.Defect(model, index, maxCharCount,
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
			if info != nil {
				info = ctx.LinesInfo()
				info.Range(func(index int, line [2]int) bool {
					if line[1] != char {
						ctx.Defect(model, index, line[0],
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
func verifyCharset(content []byte, encoding string) (string, bool) {

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
