package handlers

import (
	"encoding/json"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
	lsp "github.com/sourcegraph/go-lsp"
	"github.com/tjgurwara99/ruby-lsp/internal/code/analysis"
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
	h.logger.Printf("root path: %s\n", initializeParams.RootPath)
	h.index = analysis.New(initializeParams.RootPath, h.logger)
	go h.index.Start()
	result := lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Kind: &TextDocumentSyncKindFull,
			},
			CompletionProvider: &lsp.CompletionOptions{},
			DefinitionProvider: true,
		},
	}
	return result, nil
}
