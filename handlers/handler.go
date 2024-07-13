package handlers

import (
	"log"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tjgurwara99/ruby-lsp/internal/code/analysis"
)

type Handler struct {
	currentlyOpenFile []byte
	logger            *log.Logger
	language          *sitter.Language
	parser            *sitter.Parser
	tree              *sitter.Tree
	index             *analysis.Index
}

func New(l *log.Logger) *Handler {
	return &Handler{
		logger: l,
	}
}
