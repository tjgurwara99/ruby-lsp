package handlers

import (
	"log"

	sitter "github.com/smacker/go-tree-sitter"
)

type Handler struct {
	currentlyOpenFile []byte
	logger            *log.Logger
	language          *sitter.Language
	parser            *sitter.Parser
}

func New(l *log.Logger) *Handler {
	return &Handler{
		logger: l,
	}
}
