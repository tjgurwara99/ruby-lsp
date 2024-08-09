package navigate

import (
	"github.com/tjgurwara99/go-ruby-prism/parser"
)

func SubClasses(node parser.Node) []parser.Node {
	var classes []parser.Node
	switch n := node.(type) {
	case *parser.ModuleNode:
		body := n.Body
		for _, child := range body.Children() {
			if isClassNode(child) {
				classes = append(classes, child)
			}
		}
	case *parser.ClassNode:
		body := n.Body
		for _, child := range body.Children() {
			if isClassNode(child) {
				classes = append(classes, child)
			}
		}
	case *parser.SingletonClassNode:
		body := n.Body
		for _, child := range body.Children() {
			if isClassNode(child) {
				classes = append(classes, child)
			}
		}
	}
	return classes
}

func Methods(node parser.Node) []parser.Node {
	var methods []parser.Node
	switch n := node.(type) {
	case *parser.ModuleNode:
		body := n.Body
		for _, child := range body.Children() {
			if isMethodNode(child) {
				methods = append(methods, child)
			}
		}
	case *parser.ClassNode:
		body := n.Body
		for _, child := range body.Children() {
			if isMethodNode(child) {
				methods = append(methods, child)
			}
		}
	case *parser.SingletonClassNode:
		body := n.Body
		for _, child := range body.Children() {
			if isMethodNode(child) {
				methods = append(methods, child)
			}
		}
	}
	return methods
}

func SubModules(node parser.Node) []parser.Node {
	var modules []parser.Node
	switch n := node.(type) {
	case *parser.ModuleNode:
		body := n.Body
		for _, child := range body.Children() {
			if isModuleNode(child) {
				modules = append(modules, child)
			}
		}
	case *parser.ClassNode:
		body := n.Body
		for _, child := range body.Children() {
			if isModuleNode(child) {
				modules = append(modules, child)
			}
		}
	case *parser.SingletonClassNode:
		body := n.Body
		for _, child := range body.Children() {
			if isModuleNode(child) {
				modules = append(modules, child)
			}
		}
	}
	return modules
}

func isMethodNode(node parser.Node) bool {
	switch node.(type) {
	case *parser.DefNode, *parser.AliasMethodNode:
		return true
	}
	return false
}

func isClassNode(node parser.Node) bool {
	switch node.(type) {
	case *parser.ClassNode, *parser.SingletonClassNode:
		return true
	}
	return false
}

func isModuleNode(node parser.Node) bool {
	switch node.(type) {
	case *parser.ModuleNode:
		return true
	}
	return false
}
