package lexer_test

import (
	"testing"

	"github.com/tjgurwara99/ruby-lsp/lexer"
	"github.com/tjgurwara99/ruby-lsp/token"
)

func TestNextOperatorsDelimiters(t *testing.T) {
	input := `=+-(){},.!/*<> def true false if else
return
end
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedPos     int
		expectedCol     int
	}{
		{token.Assign, "=", 1, 0, 0},
		{token.Plus, "+", 1, 1, 1},
		{token.Minus, "-", 1, 2, 2},
		{token.LeftParen, "(", 1, 3, 3},
		{token.RightParen, ")", 1, 4, 4},
		{token.LeftBrace, "{", 1, 5, 5},
		{token.RightBrace, "}", 1, 6, 6},
		{token.Comma, ",", 1, 7, 7},
		{token.Dot, ".", 1, 8, 8},
		{token.Bang, "!", 1, 9, 9},
		{token.Slash, "/", 1, 10, 10},
		{token.Asterisk, "*", 1, 11, 11},
		{token.LessThan, "<", 1, 12, 12},
		{token.GreaterThan, ">", 1, 13, 13},
		{token.Def, "def", 1, 15, 15},
		{token.True, "true", 1, 19, 19},
		{token.False, "false", 1, 24, 24},
		{token.If, "if", 1, 30, 30},
		{token.Else, "else", 1, 33, 33},
		{token.Return, "return", 2, 38, 0},
		{token.End, "end", 3, 45, 0},
		{token.EOF, "", 3, 50, 5},
	}

	lexed := lexer.Lex(input)
	index := 0
	for token := range lexed.Tokens {
		if token.Type != tests[index].expectedType {
			t.Fatalf("expected token type %q, got %q for token %q", tests[index].expectedType, token.Type, token.Type)
		}
		if token.Literal != tests[index].expectedLiteral {
			t.Fatalf("expected token literal %q, got %q for token %q", tests[index].expectedLiteral, token.Literal, token.Type)
		}
		if token.Line != tests[index].expectedLine {
			t.Fatalf("expected token line %d, got %d for token %q", tests[index].expectedLine, token.Line, token.Type)
		}
		if token.Pos != tests[index].expectedPos {
			t.Fatalf("expected token pos %d, got %d for token %q", tests[index].expectedPos, token.Pos, token.Type)
		}
		if token.Col != tests[index].expectedCol {
			t.Fatalf("expected token col %d, got %d for token %q", tests[index].expectedCol, token.Col, token.Type)
		}
		index++
	}

}
