package handlers

import (
	"encoding/json"
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/sourcegraph/go-lsp"
	"github.com/tjgurwara99/ruby-lsp/internal/code/analysis"
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
	var ranges []*analysis.Range
	switch selected.Type() {
	case "constant":
		parent := selected.Parent()
		for parent != nil && (parent.Type() != "class" || parent.Type() != "module" || parent.Type() != "program") {
			parent = parent.Parent()
		}
		parentName := ""
		if parent != nil && parent.Type() != "program" {
			parentName = parent.NamedChild(0).Content(h.currentlyOpenFile)
		}
		nodes, ok := h.index.LookupConstant(selected.Content(h.currentlyOpenFile), parentName)
		if !ok {
			return nil, errors.New("no ranges found")
		}
		for _, n := range nodes {
			ranges = append(ranges, n.Range())
		}
	case "identifier":
		nodes, ok := h.index.LookupIdentifier(selected.Content(h.currentlyOpenFile))
		if !ok {
			return nil, errors.New("no ranges found")
		}
		for _, n := range nodes {
			ranges = append(ranges, n.Range())
		}
	default:
		return nil, errors.New("unknown node")
	}
	return Map(ranges, func(r *analysis.Range) lsp.Location {
		return lsp.Location{
			URI: lsp.DocumentURI(r.End.File),
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
