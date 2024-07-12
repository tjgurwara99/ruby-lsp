package handlers

import (
	"encoding/json"
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/sourcegraph/go-lsp"
	"github.com/tjgurwara99/ruby-lsp/code"
)

type DefinitionParams struct {
	/**
	 * The text document.
	 */
	TextDocument lsp.TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The position inside the text document.
	 */
	Position lsp.Position `json:"position"`
}

func (h *Handler) GoToDef(params json.RawMessage) (any, error) {
	var defParams DefinitionParams
	if err := json.Unmarshal(params, &defParams); err != nil {
		return nil, err
	}
	point := sitter.Point{
		Row:    uint32(defParams.Position.Line),
		Column: uint32(defParams.Position.Character),
	}
	selected := h.tree.RootNode().NamedDescendantForPointRange(point, point)
	if selected == nil {
		return nil, errors.New("failed")
	}
	var ranges []*code.Range
	var ok bool
	switch selected.Type() {
	case "constant":
		ranges, ok = h.index.LookupConstant(selected.Content(h.currentlyOpenFile))
		if !ok {
			return nil, errors.New("no ranges found")
		}
	case "identifier":
		h.logger.Println("identifier lookup started")
		ranges, ok = h.index.LookupIdentifier(selected.Content(h.currentlyOpenFile))
		h.logger.Println("identifier lookup finished")
		if !ok {
			h.logger.Println("identifier lookup errored")
			return nil, errors.New("no ranges found")
		}
	default:
		return nil, errors.New("unknown node")
	}
	return Map(ranges, func(r *code.Range) lsp.Location {
		return lsp.Location{
			URI: lsp.DocumentURI(r.End.FileURI),
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      r.Start.Line,
					Character: r.Start.Character,
				},
				End: lsp.Position{
					Line:      r.End.Line,
					Character: r.End.Character,
				},
			},
		}
	}), nil
}
