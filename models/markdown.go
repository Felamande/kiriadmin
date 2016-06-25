package models

import (
	"github.com/russross/blackfriday"
)

type Markdown struct {
	Content string `json:"content"`
	Xsrf    string `json:"_xsrf"`
}
type Article struct {
	FileName string   `json:"file_name"`
	Title    string   `json:"title"`
	Mtime    string   `json:"mtime"`
	Tags     []string `json:"tags"`
	Markdown
}

func (md Markdown) Convert() []byte {
	return blackfriday.MarkdownCommon([]byte(md.Content))
}
