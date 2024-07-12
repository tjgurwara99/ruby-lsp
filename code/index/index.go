package index

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/tjgurwara99/ruby-lsp/code"
)

type Index struct {
	Root    string
	Indexed bool
	Modules []*code.Module
	Classes []*code.Class
	Methods []*code.Method
}

func New(path string) *Index {
	return &Index{
		Root: path,
	}
}

func (i *Index) Start(logger *log.Logger) {
	language := ruby.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(language)
	logger.Println("started indexing")
	err := filepath.Walk(i.Root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules" || info.Name() == "npm-workspaces" || info.Name() == "vendor") {
			return filepath.SkipDir
		}
		if err == nil && strings.HasSuffix(info.Name(), ".rb") {
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			tree, err := parser.ParseCtx(context.Background(), nil, src)
			if err != nil {
				return err
			}
			i.indexFile(tree, src, path)
		}
		return nil
	})
	i.Indexed = true
	if err != nil {
		logger.Fatal("indexing failed")
	}
	logger.Println("indexing finished")
}

func (i *Index) indexFile(tree *sitter.Tree, src []byte, filepath string) error {
	iterator := sitter.NewNamedIterator(tree.RootNode(), sitter.DFSMode)
	iterator.ForEach(func(n *sitter.Node) error {
		switch n.Type() {
		case "module":
			return i.indexModule(n, src, nil, filepath)
		case "class":
			return i.indexClass(n, src, nil, filepath)
		case "method":
			return i.indexMethod(n, src, nil, filepath)
		}
		return nil
	})
	return nil
}

func (i *Index) indexModule(node *sitter.Node, src []byte, parent *sitter.Node, filepath string) error {
	var module code.Module
	name := node.NamedChild(0).Content(src)
	module.Name = name
	module.Locations = append(module.Locations, rangeFromNode(node, filepath))
	i.Modules = append(i.Modules, &module)
	iterator := sitter.NewNamedIterator(node, sitter.DFSMode)
	iterator.ForEach(func(n *sitter.Node) error {
		if node == n {
			return nil
		}
		switch n.Type() {
		case "module":
			return i.indexModule(n, src, parent, filepath)
		case "class":
			return i.indexClass(n, src, parent, filepath)
		case "method":
			return i.indexMethod(n, src, parent, filepath)
		}
		return nil
	})
	return nil
}

func rangeFromNode(node *sitter.Node, filepath string) code.Range {
	return code.Range{
		Start: code.Location{
			Line:      int(node.StartPoint().Row),
			Character: int(node.StartPoint().Column),
			FileURI:   filepath,
		},
		End: code.Location{
			Line:      int(node.EndPoint().Row),
			Character: int(node.EndPoint().Column),
			FileURI:   filepath,
		},
	}
}

func (i *Index) indexClass(node *sitter.Node, src []byte, parent *sitter.Node, filepath string) error {
	name := node.NamedChild(0).Content(src)
	class := code.Class{
		Name:  name,
		Range: rangeFromNode(node, filepath),
	}
	i.Classes = append(i.Classes, &class)
	return nil
}

func (i *Index) indexMethod(node *sitter.Node, src []byte, parent *sitter.Node, filepath string) error {
	// var m code.Method
	name := node.NamedChild(0).Content(src)
	// there is a possiblity that these are not args but the first statement of the method
	method := code.Method{
		Name:  name,
		Range: rangeFromNode(node, filepath),
	}
	i.Methods = append(i.Methods, &method)
	return nil
}

func (i *Index) LookupConstant(name string) ([]*code.Range, bool) {
	if !i.Indexed {
		return nil, false
	}
	var res []*code.Range
	for _, c := range i.Classes {
		if c.Name == name {
			res = append(res, &c.Range)
		}
	}
	for _, m := range i.Modules {
		if m.Name == name {
			res = append(res, &m.Locations[0])
		}
	}
	return res, true
}

func (i *Index) LookupIdentifier(name string) ([]*code.Range, bool) {
	if !i.Indexed {
		return nil, false
	}
	var res []*code.Range
	for _, m := range i.Methods {
		if m.Name == name {
			res = append(res, &m.Range)
		}
	}
	return res, true
}
