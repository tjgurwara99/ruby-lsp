package handlers

import (
	"context"
	"strings"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

func TestCompletion(t *testing.T) {
	language := ruby.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(language)
	source := `
`
	tree, err := sitter.ParseCtx(context.Background(), []byte(source), ruby.GetLanguage())
	if err != nil {
		t.Fatal(err)
	}
	selectedNode := tree.NamedDescendantForPointRange(sitter.Point{
		Row:    648,
		Column: 9,
	}, sitter.Point{
		Row:    648,
		Column: 9,
	})
	// var visit func(n *sitter.Node, name string, depth int)
	// visit = func(n *sitter.Node, name string, depth int) {
	// }
	var parent, currentNode *sitter.Node
	currentNode = selectedNode
	for {
		parent = currentNode.Parent()
		if parent.Type() == "method" {
			break
		}
		currentNode = parent
	}
	data := parent.Content([]byte(source))
	_ = data
}

func printNode(t *testing.T, n *sitter.Node, depth int, name string, source []byte) {
	prefix := ""
	if name != "" {
		prefix = name + ": "
	}
	t.Logf("%s%s%s [%d-%d] %s\n", strings.Repeat("    ", depth), prefix, n.Type(), n.StartByte(), n.EndByte(), source[n.StartByte():n.EndByte()])
}
