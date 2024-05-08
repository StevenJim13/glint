package c

import "github.com/stkali/glint/models"

var FileBasic = models.Model{
	Name: "FileBasic",
	Tags: []string{"basic"},
	Options: map[string]any{
		"charset": "utf-8",
		"maxLines": 1200,
		"maxLineLength": 120,
		"newLine": "\\n",
	},
	ModelFunc: func(model *models.Model, ctx models.Context) {

	},
}
