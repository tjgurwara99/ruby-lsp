package index

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tjgurwara99/go-ruby-prism/parser"
)

type Index struct {
	Root        string
	Indexed     bool
	ClassDecls  []ClassDecl
	ModuleDecls []ModuleDecl
	MethodDecls []MethodDecl
}

func New(path string) *Index {
	return &Index{
		Root: path,
	}
}

func (i *Index) Start(logger *log.Logger) error {
	p, err := parser.NewParser(context.Background())
	if err != nil {
		return err
	}
	logger.Println("started indexing")
	err = filepath.Walk(i.Root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") ||
			info.Name() == "node_modules" || info.Name() == "npm-workspaces" ||
			info.Name() == "vendor") {
			return filepath.SkipDir
		}
		if err == nil && strings.HasSuffix(info.Name(), ".rb") {
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			result, err := p.Parse(context.Background(), src)
			if err != nil {
				return err
			}
			i.indexProgram(result.Value, path, src)
		}
		return nil
	})
	if err != nil {
		logger.Fatal("indexing failed")
	}
	logger.Println("indexing finished")
	i.Indexed = true
	return nil
}

func (i *Index) indexProgram(node parser.Node, path string, src []byte) error {
	if node == nil {
		return nil
	}
	for _, child := range node.Children() {
		switch n := child.(type) {
		case *parser.ModuleNode:
			err := i.indexModule(n, path, src)
			if err != nil {
				return err
			}
		case *parser.ClassNode:
			err := i.indexClass(n, path, src)
			if err != nil {
				return err
			}
		case *parser.DefNode:
			err := i.indexMethod(n, path, src)
			if err != nil {
				return err
			}
		case *parser.StatementsNode:
			err := i.indexProgram(n, path, src)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func locationFromOffset(src []byte, path string, offset int) (*Location, error) {
	if offset < 0 || offset >= len(src) {
		return nil, fmt.Errorf("fileOffset is out of bounds")
	}

	line := 0
	lineOffset := 0

	for i := 0; i < offset; i++ {
		if src[i] == '\n' {
			line++
			lineOffset = 0
		} else {
			lineOffset++
		}
	}
	return &Location{
		Line:      line,
		Character: lineOffset,
		FileURI:   path,
	}, nil
}

func (i *Index) indexModule(node *parser.ModuleNode, path string, src []byte) error {
	startLocation, err := locationFromOffset(src, path, int(node.Modulekeywordloc.StartOffset))
	if err != nil {
		return err
	}
	endLocation, err := locationFromOffset(src, path, int(node.Endkeywordloc.EndOffset()))
	if err != nil {
		return err
	}
	module := ModuleDecl{
		Name: node.Name,
		r: &Range{
			Start: startLocation,
			End:   endLocation,
		},
	}
	i.ModuleDecls = append(i.ModuleDecls, module)
	return i.indexProgram(node.Body, path, src)
}

func (i *Index) indexClass(node *parser.ClassNode, path string, src []byte) error {
	startLoc, err := locationFromOffset(src, path, int(node.Classkeywordloc.StartOffset))
	if err != nil {
		return err
	}
	endLoc, err := locationFromOffset(src, path, int(node.Endkeywordloc.EndOffset()))
	if err != nil {
		return err
	}
	cls := ClassDecl{
		Name: node.Name,
		r: &Range{
			Start: startLoc,
			End:   endLoc,
		},
	}
	i.ClassDecls = append(i.ClassDecls, cls)
	return i.indexProgram(node.Body, path, src)
}

func (i *Index) indexMethod(node *parser.DefNode, path string, src []byte) error {
	startLoc, err := locationFromOffset(src, path, int(node.Defkeywordloc.StartOffset))
	if err != nil {
		return err
	}
	endLoc := startLoc
	if node.Endkeywordloc != nil {
		endLoc, err = locationFromOffset(src, path, int(node.Endkeywordloc.EndOffset()))
		if err != nil {
			return err
		}
	}
	method := MethodDecl{
		Name: node.Name,
		r: &Range{
			Start: startLoc,
			End:   endLoc,
		},
	}
	i.MethodDecls = append(i.MethodDecls, method)
	return nil
}

func (i *Index) LookupConstant(constant string) ([]*Range, bool) {
	if !i.Indexed {
		return nil, false
	}
	classRanges := mapp(filter(i.ClassDecls, func(c ClassDecl) bool {
		return c.Name == constant
	}), func(c ClassDecl) *Range {
		return c.Range()
	})
	moduleRanges := mapp(filter(i.ModuleDecls, func(c ModuleDecl) bool {
		return c.Name == constant
	}), func(c ModuleDecl) *Range {
		return c.Range()
	})
	var res []*Range
	res = append(res, moduleRanges...)
	res = append(res, classRanges...)
	return res, len(res) > 0
}

func (i *Index) LookupIdentifier(ident string) ([]*Range, bool) {
	if !i.Indexed {
		return nil, false
	}
	ranges := mapp(filter(i.MethodDecls, func(c MethodDecl) bool {
		return c.Name == ident
	}), func(c MethodDecl) *Range {
		return c.Range()
	})
	return ranges, len(ranges) > 0
}

func filter[S ~[]E, E any](s S, f func(E) bool) S {
	var res S
	for _, v := range s {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func mapp[S ~[]E, E any, R any](s S, f func(E) R) []R {
	var res []R
	for _, v := range s {
		res = append(res, f(v))
	}
	return res
}
