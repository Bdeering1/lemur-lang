package lexer

import (
    "testing"

    "lemur/token"
)

func TestNextToken(t *testing.T) {
    var input = `
        -!*/<>==!=;
        let add = fn(x, y) {
            return x + y
        }
        add(5, 10)

        if 5 < 10 { true } else { false }

        "foo"
        "foo bar"
    `
    tests := []token.Token{
        createToken("-"),
        createToken("!"),
        createToken("*"),
        createToken("/"),
        createToken("<"),
        createToken(">"),
        createToken("=="),
        createToken("!="),
        createToken(";"),

        createToken("let"),
        createIdent("add"),
        createToken("="),
        createToken("fn"),
        createToken("("),
        createIdent("x"),
        createToken(","),
        createIdent("y"),
        createToken(")"),
        createToken("{"),
        createToken("return"),
        createIdent("x"),
        createToken("+"),
        createIdent("y"),
        createToken("}"),

        createIdent("add"),
        createToken("("),
        createInt("5"),
        createToken(","),
        createInt("10"),
        createToken(")"),

        createToken("if"),
        createInt("5"),
        createToken("<"),
        createInt("10"),
        createToken("{"),
        createToken("true"),
        createToken("}"),
        createToken("else"),
        createToken("{"),
        createToken("false"),
        createToken("}"),

        createString("foo"),
        createString("foo bar"),

        createToken("\x00"),
    }

    l := New(input)

    for i, tt := range tests {
        tok := l.NextToken()
        if tok.Type != tt.Type {
            t.Fatalf("token %d: token type wrong. Expected %q, got %q",
                i + 1, tt.Type, tok.Type)
        }
        if tok.Literal != tt.Literal {
            t.Fatalf("token %d: token literal wrong. Expected %q, got %q",
                i, tt.Literal, tok.Literal)
        }
    }
}

func createIdent(l string) token.Token {
    return token.Token{Type: token.Ident, Literal: l}
}

func createInt(l string) token.Token {
    return token.Token{Type: token.Int, Literal: l }
}

func createString(l string) token.Token {
    return token.Token{Type: token.String, Literal: l }
}

func createToken(l string) (t token.Token) {
    t.Literal = l;
    switch l {

    case "\x00": t.Type = token.EOF
    case ",": t.Type = token.Comma
    case ";": t.Type = token.Semicolon
    case "(": t.Type = token.LParen
    case ")": t.Type = token.RParen
    case "{": t.Type = token.LBrace
    case "}": t.Type = token.RBrace
    case "=": t.Type = token.Assign
    case "+": t.Type = token.Plus
    case "-": t.Type = token.Minus
    case "!": t.Type = token.Bang
    case "*": t.Type = token.Asterisk
    case "/": t.Type = token.Slash
    case "<": t.Type = token.LT
    case ">": t.Type = token.GT
    case "==": t.Type = token.Eq
    case "!=": t.Type = token.NotEq
    case "fn": t.Type = token.Function
    case "let": t.Type = token.Let
    case "true": t.Type = token.True
    case "false": t.Type = token.False
    case "if": t.Type = token.If
    case "else": t.Type = token.Else
    case "return": t.Type = token.Return
    }

    return
}
