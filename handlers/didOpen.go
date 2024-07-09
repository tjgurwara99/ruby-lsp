package handlers

import (
	"encoding/json"

	"github.com/sourcegraph/go-lsp"
)

func (h *Handler) DidOpenHandler(params json.RawMessage) error {
	var paramsData lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(params, &paramsData); err != nil {
		return err
	}
	h.currentlyOpenFile = []byte(paramsData.TextDocument.Text)
	return nil
}
