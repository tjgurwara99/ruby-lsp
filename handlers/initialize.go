package handlers

import (
	"encoding/json"

	"github.com/sourcegraph/go-lsp"
)

var TextDocumentSyncKindFull = lsp.TDSKFull

func (h *Handler) Initialize(params json.RawMessage) (any, error) {
	var initializeParams lsp.InitializeParams
	if err := json.Unmarshal(params, &initializeParams); err != nil {
		return nil, err
	}
	h.logger.Println("Initialised the lsp")
	result := lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Kind: &TextDocumentSyncKindFull,
			},
			CompletionProvider: &lsp.CompletionOptions{},
		},
	}
	h.logger.Printf("%+v\n", initializeParams)
	h.logger.Printf("%+v\n", result)
	return result, nil
}
