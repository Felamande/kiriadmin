package models

import (
	"github.com/russross/blackfriday"
)

type Markdown struct {
	Content string `json:"content"  binding:"Required"`
	Xsrf    string `json:"_xsrf"  binding:"Required"`
}
type Article struct {
	FileName string   `json:"file_name" binding:"Required"`
	Title    string   `json:"title" binding:"Required"`
	Mtime    string   `json:"mtime" binding:"Required"`
	Tags     []string `json:"tags" binding:"Required"`
	Markdown
}

func (md Markdown) Convert() []byte {
	return blackfriday.MarkdownCommon([]byte(md.Content))
}
