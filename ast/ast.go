package ast

import "github.com/tjgurwara99/ruby-lsp/token"

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

type AssignStatement struct {
	Token token.Token
	Name  *Ident
	Value Expression
}

type Ident struct {
	Token token.Token
}
