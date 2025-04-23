package lexer

import (
    "testing"

    "monkey/token"
)

func TestNextToken(t *testing.T) {
    input := "=+;,(){}"

    tests := []struct{
        expectedType    token.TokenType
        expectedLiteral string
    }{
        {token.Assign, "="},
        {token.Plus, "+"},
        {token.Comma, ","},
        {token.Semicolon, ";"},
        {token.LParen, "("},
        {token.RParen, ")"},
        {token.LBrace, "{"},
        {token.RBrace, "}"},
        {token.EOF, ""},
    }

    l := New(input)

    for i, tt := range tests {
        tok := l.NextToken()
        if tok.Type != tt.expectedType {
            t.Fatalf("tests[%d] - token type wrong. Expected %q, got %q",
                i, tt.expectedType, tok.Type)
        }
        if tok.Literal != tt.expectedLiteral {
            t.Fatalf("tests[%d] - token literal wrong. Expected %q, got %q",
                i, tt.expectedType, tok.Type)
        }
    }
}
