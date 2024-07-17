package handlers

import (
	"encoding/json"
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/sourcegraph/go-lsp"
	"github.com/tjgurwara99/ruby-lsp/code/index"
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
	doc, ok := h.files[string(defParams.TextDocument.URI)]
	if !ok {
		return nil, errors.New("unopened file")
	}
	selected := doc.tree.RootNode().NamedDescendantForPointRange(point, point)
	if selected == nil {
		return nil, errors.New("failed")
	}
	var ranges []*index.Range
	nesting := findNesting(doc.content, selected)
	switch selected.Type() {
	case "constant":
		ranges, ok = h.index.LookupConstant(selected.Content(doc.content), nesting)
		if !ok {
			return nil, errors.New("no ranges found")
		}
	case "identifier":
		ranges, ok = h.index.LookupIdentifier(selected.Content(doc.content), nesting)
		if !ok {
			h.logger.Println("identifier lookup errored")
			return nil, errors.New("no ranges found")
		}
	default:
		h.logger.Printf("unknown node type %s", selected.Type())
		return nil, errors.New("unknown node")
	}
	return Map(ranges, func(r *index.Range) lsp.Location {
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
