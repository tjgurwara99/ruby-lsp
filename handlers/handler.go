package handlers

import (
	"log"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tjgurwara99/ruby-lsp/code/index"
)

type Handler struct {
	currentlyOpenFile []byte
	logger            *log.Logger
	language          *sitter.Language
	parser            *sitter.Parser
	tree              *sitter.Tree
	index             *index.Index
}

func New(l *log.Logger) *Handler {
	return &Handler{
		logger: l,
	}
}
