package handlers

import (
	"context"
	"encoding/json"
	"errors"

	lsp "github.com/sourcegraph/go-lsp"
)

func (h *Handler) DidOpenHandler(params json.RawMessage) error {
	var paramsData lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(params, &paramsData); err != nil {
		return err
	}
	tree, err := h.parser.ParseCtx(context.Background(), nil, []byte(paramsData.TextDocument.Text))
	if err != nil {
		return err
	}
	h.files[string(paramsData.TextDocument.URI)] = &TextDocument{
		content: []byte(paramsData.TextDocument.Text),
		tree:    tree,
	}
	return nil
}

func (h *Handler) DidChangeHandler(params json.RawMessage) error {
	var paramsData lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(params, &paramsData); err != nil {
		return err
	}
	doc, ok := h.files[string(paramsData.TextDocument.URI)]
	if !ok {
		return errors.New("file never opened")
	}
	doc.content = []byte(paramsData.ContentChanges[0].Text)

	var err error
	doc.tree, err = h.parser.ParseCtx(context.Background(), doc.tree, doc.content)
	if err != nil {
		return err
	}
	return nil
}
