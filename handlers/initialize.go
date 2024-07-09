package handlers

import (
	"encoding/json"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/sourcegraph/go-lsp"
)

var TextDocumentSyncKindFull = lsp.TDSKFull

func (h *Handler) Initialize(params json.RawMessage) (any, error) {
	var initializeParams lsp.InitializeParams
	if err := json.Unmarshal(params, &initializeParams); err != nil {
		return nil, err
	}
	h.language = ruby.GetLanguage()
	h.parser = sitter.NewParser()
	h.parser.SetLanguage(h.language)
	result := lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Kind: &TextDocumentSyncKindFull,
			},
			CompletionProvider: &lsp.CompletionOptions{},
		},
	}
	return result, nil
}
