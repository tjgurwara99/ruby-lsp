package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/stretchr/testify/require"
)

func TestNesting(t *testing.T) {
	language := ruby.GetLanguage()
	p := sitter.NewParser()
	p.SetLanguage(language)
	f, err := os.Open("../code/index/testdata/foo.rb")
	require.NoError(t, err)
	code, err := io.ReadAll(f)
	require.NoError(t, err)
	tree, err := p.ParseCtx(context.Background(), nil, code)
	require.NoError(t, err)
	point := sitter.Point{
		Row:    8,
		Column: 16,
	}
	node := tree.RootNode().NamedDescendantForPointRange(point, point)
	nesting := findNesting(code, node)
	fmt.Println(nesting)
}
