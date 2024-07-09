package handlers

import (
	"encoding/json"

	"github.com/sourcegraph/go-lsp"
)

func (h *Handler) TextCompletion(params json.RawMessage) (any, error) {
	var completion lsp.CompletionParams
	if err := json.Unmarshal(params, &completion); err != nil {
		return nil, err
	}
	// find all possible things
	h.logger.Printf("Text Completion: %+v", completion)
	return lsp.CompletionItem{}, nil
}
