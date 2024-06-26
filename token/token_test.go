package token_test

import (
	"testing"

	"github.com/tjgurwara99/ruby-lsp/token"
)

func TestIsLiteral(t *testing.T) {
	testCases := []struct {
		tokenType token.TokenType
		expected  bool
	}{
		{token.Illegal, false},
		{token.EOF, false},
		{token.Ident, true},
		{token.Int, true},
		{token.Assign, false},
		{token.Plus, false},
		{token.Minus, false},
		{token.Bang, false},
		{token.Asterisk, false},
		{token.Slash, false},
		{token.LessThan, false},
		{token.GreaterThan, false},
		{token.Comma, false},
		{token.LeftParen, false},
		{token.RightParen, false},
		{token.LeftBrace, false},
		{token.RightBrace, false},
		{token.And, false},
		{token.Or, false},
		{token.Eq, false},
		{token.NotEq, false},
		{token.LessThanEq, false},
		{token.GreaterThanEq, false},
		{token.Class, false},
		{token.Module, false},
		{token.Def, false},
		{token.Do, false},
		{token.True, false},
		{token.False, false},
		{token.If, false},
		{token.End, false},
		{token.Return, false},
	}

	for _, tc := range testCases {
		if tc.expected != tc.tokenType.IsLiteral() {
			t.Errorf("%s is literal: %t", tc.tokenType, tc.expected)
		}
	}
}

func TestIsOperator(t *testing.T) {
	testCases := []struct {
		tokenType token.TokenType
		expected  bool
	}{
		{token.Illegal, false},
		{token.EOF, false},
		{token.Ident, false},
		{token.Int, false},
		{token.Class, false},
		{token.Module, false},
		{token.Def, false},
		{token.Do, false},
		{token.True, false},
		{token.False, false},
		{token.If, false},
		{token.End, false},
		{token.Return, false},
	}

	for _, tc := range testCases {
		if tc.expected != tc.tokenType.IsOperator() {
			t.Errorf("%s is literal: %t", tc.tokenType, tc.expected)
		}
	}
}

func TestTokenType_IsKeyword(t *testing.T) {
	testCases := []struct {
		tokenType token.TokenType
		expected  bool
	}{
		{token.Illegal, false},
		{token.EOF, false},
		{token.Ident, false},
		{token.Int, false},
		{token.Class, true},
		{token.Module, true},
		{token.Def, true},
		{token.Do, true},
		{token.True, true},
		{token.False, true},
		{token.If, true},
		{token.End, true},
		{token.Return, true},
		{token.Illegal, false},
		{token.EOF, false},
		{token.Ident, false},
		{token.Int, false},
		{token.Assign, false},
		{token.Plus, false},
		{token.Minus, false},
		{token.Bang, false},
		{token.Asterisk, false},
		{token.Slash, false},
		{token.LessThan, false},
		{token.GreaterThan, false},
		{token.Comma, false},
		{token.LeftParen, false},
		{token.RightParen, false},
		{token.LeftBrace, false},
		{token.RightBrace, false},
		{token.And, false},
		{token.Or, false},
		{token.Eq, false},
		{token.NotEq, false},
		{token.LessThanEq, false},
		{token.GreaterThanEq, false},
	}

	for _, tc := range testCases {
		if tc.expected != tc.tokenType.IsKeyword() {
			t.Errorf("%s is literal: %t", tc.tokenType, tc.expected)
		}
	}
}

func TestIsKeyword(t *testing.T) {
	testCases := []struct {
		keyword  string
		expected bool
	}{
		{"", false},
		{"a", false},
		{"def", true},
		{"true", true},
		{"false", true},
		{"if", true},
		{"else", true},
		{"return", true},
	}

	for _, tc := range testCases {
		if tc.expected != token.IsKeyword(tc.keyword) {
			t.Errorf("%s is keyword: %t", tc.keyword, tc.expected)
		}
	}
}

func TestIsIdentifier(t *testing.T) {
	testCases := []struct {
		identifier string
		expected   bool
	}{
		{"", false},
		{"a", true},
		{"def", false},
		{"true", false},
		{"false", false},
		{"if", false},
		{"else", false},
		{"return", false},
		{"_abc", true},
		{"a-bc", false},
	}

	for _, tc := range testCases {
		if tc.expected != token.IsIdentifier(tc.identifier) {
			t.Errorf("%s is identifier: %t", tc.identifier, tc.expected)
		}
	}
}
