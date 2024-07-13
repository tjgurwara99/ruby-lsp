package handlers

import (
	"log"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tjgurwara99/ruby-lsp/code/index"
)

type TextDocument struct {
	tree    *sitter.Tree
	content []byte
}

type Handler struct {
	logger   *log.Logger
	language *sitter.Language
	parser   *sitter.Parser
	files    map[string]*TextDocument
	index    *index.Index
}

func New(l *log.Logger) *Handler {
	return &Handler{
		logger: l,
		files:  make(map[string]*TextDocument),
	}
}
