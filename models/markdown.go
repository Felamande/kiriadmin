package models

import (
	"github.com/russross/blackfriday"
)

type Markdown struct {
	Content string `form:"content"`
}

func (md Markdown) Eval() []byte {
	return blackfriday.MarkdownCommon([]byte(md.Content))
}
