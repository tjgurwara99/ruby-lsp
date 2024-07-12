package handlers

import (
	"context"
	"encoding/json"

	lsp "github.com/sourcegraph/go-lsp"
)

func (h *Handler) DidOpenHandler(params json.RawMessage) error {
	var paramsData lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(params, &paramsData); err != nil {
		return err
	}
	h.currentlyOpenFile = []byte(paramsData.TextDocument.Text)
	var err error
	h.tree, err = h.parser.ParseCtx(context.Background(), nil, h.currentlyOpenFile)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DidChangeHandler(params json.RawMessage) error {
	var paramsData lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(params, &paramsData); err != nil {
		return err
	}
	h.currentlyOpenFile = []byte(paramsData.ContentChanges[0].Text)
	var err error
	h.tree, err = h.parser.ParseCtx(context.Background(), h.tree, h.currentlyOpenFile)
	if err != nil {
		return err
	}
	return nil
}
