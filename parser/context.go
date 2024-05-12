package parser

import (
	"fmt"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/tool"
)

type ASTContext struct {
	file    string
	info    LinesInfo
	content []byte
	ast     *sitter.Tree
}

// Name implements Context.
func (c *ASTContext) Name() string {
	return c.file
}

// AddDefect implements parser.Context.
func (c *ASTContext) AddDefect(defect Defect) {
	fmt.Println("-------", defect)
	return
}

// CallExpresses implements parser.Context.
func (c *ASTContext) CallExpresses() []*CallExpress {
	return nil
}

// Content implements parser.Context.
func (c *ASTContext) Content() []byte {
	if c.content == nil {
		if err := c.loadContent(); err != nil {
			errors.Warningf("failed to get file: %q content, err: %s", c.file, err)
		}
	}
	return c.content
}

// loadContent TODO
func (c *ASTContext) loadContent() (err error) {
	c.content, err = os.ReadFile(c.file)
	return
}

// Functions implements parser.Context.
func (c *ASTContext) Functions() []*Function {
	return nil
}

// LinesInfo implements parser.Context.
func (c *ASTContext) LinesInfo() LinesInfo {
	if c.info == nil {
		c.getLinesInfo()
	}
	return c.info
}

// -------- ------------0
// mac 		\r			1
// linux 	\n			2
// windows 	\r\n		3
func (c *ASTContext) getLinesInfo() LinesInfo {
	ctt := c.Content()
	gap := 0
	index, length := 0, len(ctt)

	for index < length {
		switch ctt[index] {
		case '\r':
			lineLength := len(tool.ToString(ctt[gap:index]))
			if index+1 < length {
				if ctt[index+1] == '\n' {
					// \r\n
					c.info = append(c.info, [2]int{lineLength, 3})
					index += 1
					gap = index + 1
				} else {
					// \r
					c.info = append(c.info, [2]int{lineLength, 1})
					gap = index + 1
				}
			} else {
				// EOF
				c.info = append(c.info, [2]int{lineLength, 1})
			}
		case '\n':
			lineLength := len(tool.ToString(ctt[gap:index]))
			c.info = append(c.info, [2]int{lineLength, 2})
			gap = index + 1
		}
		index += 1
	}
	if index > gap {
		lineLength := len(tool.ToString(ctt[gap:index]))
		c.info = append(c.info, [2]int{lineLength, 2})
	}
	return nil
}

var _ Context = (*ASTContext)(nil)
