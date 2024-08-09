package index

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/tjgurwara99/go-ruby-prism/parser"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	"testing"
)

func TestIndex(t *testing.T) {
	i := New("/Users/taj/github/github")
	i.Start(log.Default())
}

func TestPrism(t *testing.T) {
	ctx := context.Background()
	p, err := parser.NewParser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("./testdata/foo.rb")
	if err != nil {
		t.Fatal(err)
	}
	result, err := p.Parse(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	dfs(t, result.Value, 0)
}

func dfs(t *testing.T, node parser.Node, indent int) {
	t.Helper()
	t.Logf("%s%T\n", strings.Repeat("  ", indent), node)
	for _, child := range node.Children() {
		dfs(t, child, indent+1)
	}
}

func BenchmarkParsing(b *testing.B) {
	ctx := context.Background()
	p, err := parser.NewParser(ctx)
	if err != nil {
		b.Fatal(err)
	}
	data, err := os.ReadFile("./testdata/bat/bat.rb")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Parse(ctx, data)
	}
}

func BenchmarkTreeSitterParser(b *testing.B) {
	language := ruby.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(language)
	src, err := os.ReadFile("./testdata/bat/bat.rb")
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseCtx(ctx, nil, src)
	}
}
