package token

import (
	"fmt"
	"unicode"
)

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Pos     int
	Col     int
}

const (
	Illegal TokenType = iota
	EOF

	literalBeginning
	Ident
	Int
	literalEnd

	operatorBeginning
	Assign
	Plus
	Minus
	Bang
	Asterisk
	Slash

	LessThan
	GreaterThan

	Comma
	SemiColon

	LeftParen
	RightParen
	LeftBrace
	RightBrace

	And
	Or
	Eq
	NotEq
	LessThanEq
	GreaterThanEq
	operatorEnd

	keywordBeginning
	Class
	Module
	Def
	Do
	True
	False
	If
	Else
	End
	Return
	keywordEnd
)

var tokens = [...]string{
	Illegal:       "Illegal",
	EOF:           "EOF",
	Ident:         "Ident",
	Int:           "int",
	Assign:        "=",
	Plus:          "+",
	Minus:         "-",
	Bang:          "!",
	Asterisk:      "*",
	Slash:         "/",
	LessThan:      "<",
	GreaterThan:   ">",
	Comma:         ",",
	SemiColon:     ";",
	LeftParen:     "(",
	RightParen:    ")",
	LeftBrace:     "{",
	RightBrace:    "}",
	And:           "&&",
	Or:            "||",
	Eq:            "==",
	NotEq:         "!=",
	LessThanEq:    "<=",
	GreaterThanEq: ">=",
	Class:         "class",
	Module:        "module",
	Def:           "def",
	Do:            "do",
	True:          "true",
	False:         "false",
	If:            "if",
	Else:          "else",
	End:           "end",
	Return:        "return",
}

func (t TokenType) String() string {
	if t == literalBeginning || t == literalEnd || t == operatorBeginning || t == operatorEnd || t == keywordBeginning || t == keywordEnd {
		return fmt.Sprintf("token(%d)", t)
	}

	if 0 < t && t < keywordEnd {
		return tokens[t]
	}
	return fmt.Sprintf("token(%d)", t)
}

var keywordsToTokenType map[string]TokenType

func init() {
	if keywordsToTokenType == nil {
		keywordsToTokenType = make(map[string]TokenType)
	}
	for i := keywordBeginning + 1; i < keywordEnd; i++ {
		keywordsToTokenType[tokens[i]] = i
	}
}

func (t TokenType) IsLiteral() bool {
	return literalBeginning < t && t < literalEnd
}

func (t TokenType) IsOperator() bool {
	return operatorBeginning < t && t < operatorEnd
}

func (t TokenType) IsKeyword() bool {
	return keywordBeginning < t && t < keywordEnd
}

func IsKeyword(ident string) bool {
	_, ok := keywordsToTokenType[ident]
	return ok
}

func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && (i == 0 || !unicode.IsDigit(c)) && c != '_' {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}

func TokenForIdent(name string) TokenType {
	if tok, ok := keywordsToTokenType[name]; ok {
		return tok
	}
	return Ident
}
