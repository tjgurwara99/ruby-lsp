package handlers

import (
	"context"
	"encoding/json"

	sitter "github.com/smacker/go-tree-sitter"
	lsp "github.com/sourcegraph/go-lsp"
)

func (h *Handler) TextCompletion(params json.RawMessage) (any, error) {
	var paramsData lsp.CompletionParams
	if err := json.Unmarshal(params, &paramsData); err != nil {
		return nil, err
	}
	tree := h.tree
	h.logger.Println("POSITION", paramsData.Position)
	allIdents, err := sitter.NewQuery([]byte(allIdentsQuery), h.language)
	if err != nil {
		return nil, err
	}
	queryCursor := sitter.NewQueryCursor()
	queryCursor.Exec(allIdents, tree.RootNode())
	for {
		m, ok := queryCursor.NextMatch()
		if !ok {
			h.logger.Println("Query match not found")
			break
		}
		h.logger.Println("match:", m)

		m = queryCursor.FilterPredicates(m, h.currentlyOpenFile)
		for _, c := range m.Captures {
			h.logger.Println("content in the node:", c.Node.Content(h.currentlyOpenFile))
		}

	}
	data, err := allIdentifiers(h.currentlyOpenFile, h.language, h.parser, paramsData)
	if err != nil {
		return nil, err
	}
	h.logger.Printf("All idents: %+v", data)
	// find all possible things
	return lsp.CompletionList{
		IsIncomplete: false,
		Items:        data,
	}, nil
}

const allClassNamesQuery = `((class
	name: [
		(constant) @clsName
    ]
))`

const allMethodNamesQuery = `((method
	name: (identifier) @ident
))`

const allIdentsQuery = `((identifier) @ident)`

const parentClassScopedQueryForIdentifiers = `((class
    body: [
    	(body_statement
        	(method
        		name: (identifier) @ident))
    ]
))`

const parentModuleScopedQueryForIdents = `((module
    body: [
    	(body_statement
        	(method
        		name: (identifier) @ident))
    ]
))`

// const methodQuery = `((
// 	(method
// 		name: [
// 			(identifier) @name
// 		]
// 		_ @params
//         _ @body
//      )
// ) @method
// (#eq? @name "some_integer"))`

func Map[T, V any](items []T, fn func(T) V) []V {
	if items == nil {
		return nil
	}
	out := make([]V, 0, len(items))
	for _, item := range items {
		out = append(out, fn(item))
	}
	return out
}

func scopedSourceContent(src []byte, lang *sitter.Language, params lsp.CompletionParams) (string, error) {
	tree, err := sitter.ParseCtx(context.Background(), []byte(src), lang)
	if err != nil {
		return "", err
	}
	selectedNode := tree.NamedDescendantForPointRange(sitter.Point{
		Row:    uint32(params.Position.Line),
		Column: uint32(params.Position.Character),
	}, sitter.Point{
		Row:    uint32(params.Position.Line),
		Column: uint32(params.Position.Character),
	})

	var parent, currentNode *sitter.Node
	currentNode = selectedNode
	for {
		parent = currentNode.Parent()
		if parent.Type() == "method" || parent.Type() == "class" || parent.Type() == "module" || parent.Type() == "program" {
			break
		}
		currentNode = parent
	}
	return parent.Content([]byte(src)), nil
}

// func allScopedIdentifiers(src []byte, lang *sitter.Language, parser *sitter.Parser, params lsp.CompletionParams) ([]lsp.CompletionItem, error) {
// 	scopedSrc, err := scopedSourceContent(src, lang, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	allScopedIdents, err := executeQuery(allIdentsQuery, []byte(scopedSrc), lang, parser)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }

func allIdentifiers(data []byte, lang *sitter.Language, parser *sitter.Parser, params lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	classNames, err := executeQuery(allClassNamesQuery, data, lang, parser)
	if err != nil {
		return nil, err
	}
	methodNames, err := executeQuery(allMethodNamesQuery, data, lang, parser)
	if err != nil {
		return nil, err
	}
	allIdents, err := executeQuery(allIdentsQuery, data, lang, parser)
	if err != nil {
		return nil, err
	}
	allIdentsMap := toSet(allIdents)
	classes := Map(classNames, func(s string) lsp.CompletionItem {
		delete(allIdentsMap, s)
		return lsp.CompletionItem{
			Label: s,
			Kind:  lsp.CIKClass,
		}
	})
	methods := Map(methodNames, func(s string) lsp.CompletionItem {
		delete(allIdentsMap, s)
		return lsp.CompletionItem{
			Label: s,
			Kind:  lsp.CIKMethod,
		}
	})
	allIdents = toSlice(allIdentsMap)
	idents := Map(allIdents, func(s string) lsp.CompletionItem {
		return lsp.CompletionItem{
			Label: s,
			Kind:  lsp.CIKVariable,
		}
	})
	return join(classes, methods, idents), nil
}

func toSet[T comparable](data []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, e := range data {
		m[e] = struct{}{}
	}
	return m
}

func join[T any](a []T, rest ...[]T) []T {
	for _, d := range rest {
		for _, v := range d {
			a = append(a, v)
		}
	}
	return a
}

func toSlice[T comparable, M any](data map[T]M) []T {
	var slice []T
	for e := range data {
		slice = append(slice, e)
	}
	return slice
}

func executeQuery(query string, src []byte, lang *sitter.Language, parser *sitter.Parser) ([]string, error) {
	tree, err := parser.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil, err
	}
	q, err := sitter.NewQuery([]byte(query), lang)
	if err != nil {
		return nil, err
	}
	queryCursor := sitter.NewQueryCursor()
	queryCursor.Exec(q, tree.RootNode())
	var idents []string
	for {
		m, ok := queryCursor.NextMatch()
		if !ok {
			break
		}

		m = queryCursor.FilterPredicates(m, src)
		for _, c := range m.Captures {
			idents = append(idents, c.Node.Content(src))
		}
	}
	return idents, nil
}
