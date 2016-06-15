package models

import (
	"github.com/russross/blackfriday"
)

type Markdown struct {
	FileName string   `json:"file_name"`
	Title    string   `json:"title"`
	Mtime    string   `json:"mtime"`
	Tags     []string `json:"tags"`
	Content  string   `json:"content"`
	Xsrf     string   `json:"_xsrf"`
}

func (md Markdown) Convert() []byte {
	return blackfriday.MarkdownCommon([]byte(md.Content))
}
