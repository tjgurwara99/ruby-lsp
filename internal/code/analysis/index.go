package analysis

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

type Index struct {
	Root    string
	Indexed bool
	logger  *log.Logger
	Classes []*ClassDecl
	Modules []*ModuleDecl
	Methods []*MethodDecl
}

func New(path string, logger *log.Logger) *Index {
	return &Index{
		Root:   path,
		logger: logger,
	}
}

func (i *Index) Start() {
	lang := ruby.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	i.logger.Println("started indexing")
	err := filepath.Walk(i.Root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || slices.Contains(ignoredDirectories, info.Name())) {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if !strings.HasSuffix(info.Name(), ".rb") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		tree, err := parser.ParseCtx(context.Background(), nil, src)
		if err != nil {
			return err
		}
		return i.indexFile(tree, src, path)
	})
	if err != nil {
		i.logger.Printf("error walking the workspace: %s", err)
	}
	i.logger.Println("indexing finished")
	i.Indexed = true
}

func (i *Index) indexFile(tree *sitter.Tree, src []byte, filepath string) error {
	node := tree.RootNode()
	var err error
	for j := 0; j < int(node.NamedChildCount()); j++ {
		n := node.NamedChild(j)
		if node == n {
			continue
		}
		switch n.Type() {
		case "module":
			err = i.indexModule(n, src, filepath, nil)
		case "class":
			err = i.indexClass(n, src, filepath, nil)
		case "method":
			err = i.indexMethod(n, src, filepath, nil)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (i *Index) indexModule(node *sitter.Node, src []byte, filepath string, parent Node) error {
	module := ModuleDecl{
		Name: node.NamedChild(0).Content(src),
		Location: &Location{
			Line:      int(node.Range().StartPoint.Row),
			Character: int(node.Range().StartPoint.Column),
			File:      filepath,
		},
		Parent: parent,
	}
	for j := 0; j < int(node.NamedChildCount()); j++ {
		n := node.NamedChild(j)
		if node == n {
			continue
		}
		switch n.Type() {
		case "module":
			i.indexModule(n, src, filepath, &module)
		case "class":
			i.indexClass(n, src, filepath, &module)
		case "method":
			i.indexMethod(n, src, filepath, &module)
		}
	}
	i.Modules = append(i.Modules, &module)
	return nil
}

func (i *Index) indexClass(node *sitter.Node, src []byte, filepath string, parent Node) error {
	class := ClassDecl{
		Name:   node.NamedChild(0).Content(src),
		Parent: parent,
		r: &Range{
			Start: &Location{
				Line:      int(node.StartPoint().Row),
				Character: int(node.StartPoint().Column),
				File:      filepath,
			},
			End: &Location{
				Line:      int(node.EndPoint().Row),
				Character: int(node.EndPoint().Column),
				File:      filepath,
			},
		},
	}
	for j := 0; j < int(node.NamedChildCount()); j++ {
		n := node.NamedChild(j)
		if node == n {
			continue
		}
		switch n.Type() {
		case "module":
			i.indexModule(n, src, filepath, &class)
		case "class":
			i.indexClass(n, src, filepath, &class)
		case "method":
			i.indexMethod(n, src, filepath, &class)
		}
	}
	i.Classes = append(i.Classes, &class)
	return nil
}

func (i *Index) indexMethod(node *sitter.Node, src []byte, filepath string, parent Node) error {
	method := MethodDecl{
		Name:   node.NamedChild(0).Content(src),
		Parent: parent,
		r: &Range{
			Start: &Location{
				Line:      int(node.StartPoint().Row),
				Character: int(node.StartPoint().Column),
				File:      filepath,
			},
			End: &Location{
				Line:      int(node.EndPoint().Row),
				Character: int(node.EndPoint().Column),
				File:      filepath,
			},
		},
	}
	i.Methods = append(i.Methods, &method)
	return nil
}

func (i *Index) LookupConstant(ident string, parentName string) ([]Node, bool) {
	if !i.Indexed {
		return nil, false
	}
	i.logger.Printf("looking up %q with parent %q", ident, parentName)
	var res []Node
	for _, c := range i.Classes {
		if c.Name == ident && c.ParentName() == parentName {
			res = append(res, c)
		}
	}
	for _, m := range i.Modules {
		if m.Name == ident && m.ParentName() == parentName {
			res = append(res, m)
		}
	}
	i.logger.Printf("found %+v", res)
	return res, true
}

func (i *Index) LookupIdentifier(ident string) ([]Node, bool) {
	if !i.Indexed {
		return nil, false
	}
	i.logger.Printf("looking up %q", ident)
	var res []Node
	for _, m := range i.Methods {
		if m.Name == ident {
			res = append(res, m)
		}
	}
	i.logger.Printf("found %+v", res)
	return res, true
}
