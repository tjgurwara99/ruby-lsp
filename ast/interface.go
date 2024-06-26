package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	declare()
}

type Expression interface {
	Node
	evaluate()
}
