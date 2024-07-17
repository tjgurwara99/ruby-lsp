package handlers

import (
	"slices"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func findNesting(src []byte, node *sitter.Node) []string {
	var namespaces []string
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		if parent.Type() == "module" || parent.Type() == "class" {
			name := parent.NamedChild(0).Content(src)
			namespaces = append(namespaces, name)
		}
	}
	slices.Reverse(namespaces)
	var res []string
	for i := len(namespaces) - 1; i >= 0; i-- {
		res = append(res, strings.Join(namespaces[:i+1], "::"))
	}
	return res
}
