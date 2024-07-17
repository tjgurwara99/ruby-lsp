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
)

type Attribute struct {
	Ident string
	Type  string
}

type Symbol struct {
	Name       string
	Type       string
	Attributes []*Attribute
	r          *Range
}

type Index struct {
	Root    string
	Indexed bool
	Symbols []*Symbol
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
	node := tree.RootNode()
	for j := 0; j < int(node.ChildCount()); j++ {
		n := node.Child(j)
		switch n.Type() {
		case "module":
			_, err := i.indexModule(n, src, filepath, "")
			if err != nil {
				return err
			}
		case "class":
			_, err := i.indexClass(n, src, filepath, "")
			if err != nil {
				return err
			}
		case "method":
			_, err := i.indexMethod(n, src, filepath, "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Index) lookupSymbol(scope string, t string, r *Range) (*Symbol, bool) {
	for _, sym := range i.Symbols {
		if sym.Name == scope && sym.Type == t {
			return sym, true
		}
	}
	return &Symbol{
		Name: scope,
		Type: t,
		r:    r,
	}, false
}

func (i *Index) indexModule(node *sitter.Node, src []byte, filepath string, scope string) (*ModuleDecl, error) {
	name := node.NamedChild(0).Content(src)
	if scope == "" {
		scope = name
	} else {
		scope = scope + "::" + name
	}
	rr := rangeFromNode(node, filepath)
	symbol, symbolIndexed := i.lookupSymbol(scope, "module", rr)
	if !symbolIndexed {
		i.Symbols = append(i.Symbols, symbol)
	}
	var module ModuleDecl
	var bodyNode *sitter.Node
	for j := 0; j < int(node.NamedChildCount()); j++ {
		n := node.NamedChild(j)
		if node == n {
			continue
		}
		if n.Type() == "body_statement" {
			bodyNode = n
		}
	}
	if bodyNode == nil {
		return nil, nil
	}
	var attributes []*Attribute
	for j := 0; j < int(bodyNode.NamedChildCount()); j++ {
		n := bodyNode.NamedChild(j)
		if n == bodyNode {
			continue
		}
		switch n.Type() {
		case "module":
			submodule, err := i.indexModule(n, src, filepath, scope)
			if err != nil {
				return nil, err
			}
			module.Modules = append(module.Modules, submodule)
		case "class":
			class, err := i.indexClass(n, src, filepath, scope)
			if err != nil {
				return nil, err
			}
			module.Classes = append(module.Classes, class)
		case "method":
			method, err := i.indexMethod(n, src, filepath, scope)
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, &Attribute{
				Ident: method.Name,
				Type:  "method",
			})
			module.Methods = append(module.Methods, method)
		}
	}
	symbol.Attributes = attributes
	return &module, nil
}

func rangeFromNode(node *sitter.Node, filepath string) *Range {
	return &Range{
		Start: &Location{
			Line:      int(node.StartPoint().Row),
			Character: int(node.StartPoint().Column),
			FileURI:   filepath,
		},
		End: &Location{
			Line:      int(node.EndPoint().Row),
			Character: int(node.EndPoint().Column),
			FileURI:   filepath,
		},
	}
}

func (i *Index) indexClass(node *sitter.Node, src []byte, filepath string, scope string) (*ClassDecl, error) {
	name := node.NamedChild(0).Content(src)
	rr := rangeFromNode(node, filepath)
	class := ClassDecl{
		Name: name,
		r:    rr,
	}
	if scope == "" {
		scope = name
	} else {
		scope = scope + "::" + name
	}
	symbol, symbolIndexed := i.lookupSymbol(scope, "module", rr)
	if !symbolIndexed {
		i.Symbols = append(i.Symbols, symbol)
	}
	var bodyNode *sitter.Node
	for j := 0; j < int(node.NamedChildCount()); j++ {
		n := node.NamedChild(j)
		if node == n {
			continue
		}
		if n.Type() == "body_statement" {
			bodyNode = n
		}
	}
	if bodyNode == nil {
		return &class, nil
	}
	var attributes []*Attribute
	for j := 0; j < int(bodyNode.NamedChildCount()); j++ {
		n := bodyNode.NamedChild(j)
		if n == bodyNode {
			continue
		}
		switch n.Type() {
		case "module":
			submodule, err := i.indexModule(n, src, filepath, scope)
			if err != nil {
				return nil, err
			}
			class.Modules = append(class.Modules, submodule)
		case "class":
			class, err := i.indexClass(n, src, filepath, scope)
			if err != nil {
				return nil, err
			}
			class.Classes = append(class.Classes, class)
		case "method":
			method, err := i.indexMethod(n, src, filepath, scope)
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, &Attribute{
				Ident: method.Name,
				Type:  "method",
			})
			class.Methods = append(class.Methods, method)
		}
	}
	symbol.Attributes = attributes
	return &class, nil
}

func (i *Index) indexMethod(node *sitter.Node, src []byte, filepath string, scope string) (*MethodDecl, error) {
	// var m Method
	name := node.NamedChild(0).Content(src)
	rr := rangeFromNode(node, filepath)
	// there is a possiblity that these are not args but the first statement of the method
	n := node.NamedChild(1)
	if scope == "" {
		scope = name
	} else {
		scope = scope + "." + name
	}
	symbol, symbolIndexed := i.lookupSymbol(scope, "method", rr)
	if !symbolIndexed {
		i.Symbols = append(i.Symbols, symbol)
	}
	var args []string
	if n != nil && n.Type() == "method_parameters" {
		s := strings.TrimPrefix(n.Content(src), "(")
		s = strings.TrimSuffix(s, ")")
		args = strings.Split(s, ",")
	}
	method := MethodDecl{
		Name: name,
		r:    rangeFromNode(node, filepath),
		Args: args,
	}
	return &method, nil
}

func (i *Index) LookupConstant(name string, nesting []string) ([]*Range, bool) {
	if !i.Indexed {
		return nil, false
	}
	var res []*Range
	// use the nesting[i] + "::" + name to get the search parameter in symbols
	for i := range nesting {
		nesting[i] = nesting[i] + "::" + name
	}
	for _, sym := range i.Symbols {
		for _, n := range nesting {
			if sym.Name == n {
				res = append(res, sym.r)
			}
		}
	}
	return res, true
}

func (i *Index) LookupIdentifier(name string, nesting []string) ([]*Range, bool) {
	if !i.Indexed {
		return nil, false
	}
	var res []*Range
	// use the nesting[i] +"."+ name to get the search parameter in symbols
	for i := range nesting {
		nesting[i] = nesting[i] + "." + name
	}
	for _, sym := range i.Symbols {
		for _, n := range nesting {
			if sym.Name == n {
				res = append(res, sym.r)
			}
		}
	}
	return res, true
}
