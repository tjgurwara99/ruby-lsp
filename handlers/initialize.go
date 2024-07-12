package handlers

import (
	"encoding/json"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
	lsp "github.com/sourcegraph/go-lsp"
	"github.com/tjgurwara99/ruby-lsp/code/index"
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
	h.index = index.New(initializeParams.RootPath)
	go h.index.Start(h.logger)
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
