package lexer

import (
    "testing"

    "lemur/token"
)

type tToken struct {
    expType    token.TokenType
    expLiteral string
}

func TestNextToken(t *testing.T) {
    var input = `
        let five = 5;
        let ten = 10;

        let add = fn(x, y) {
            x + y;
        };
        let res = add(five, ten);
        -!*/<>;
        if 5 < 10 {
            return true;
        } else {
            return false;
        }
        5 == 5;
        5 != 10;
    `
    tests := []tToken{
        tTok("let"),
        tUser("five"),
        tTok("="),
        tInt("5"),
        tTok(";"),
        tTok("let"),
        tUser("ten"),
        tTok("="),
        tInt("10"),
        tTok(";"),
        tTok("let"),
        tUser("add"),
        tTok("="),
        tTok("fn"),
        tTok("("),
        tUser("x"),
        tTok(","),
        tUser("y"),
        tTok(")"),
        tTok("{"),
        tUser("x"),
        tTok("+"),
        tUser("y"),
        tTok(";"),
        tTok("}"),
        tTok(";"),
        tTok("let"),
        tUser("res"),
        tTok("="),
        tUser("add"),
        tTok("("),
        tUser("five"),
        tTok(","),
        tUser("ten"),
        tTok(")"),
        tTok(";"),
        tTok("-"),
        tTok("!"),
        tTok("*"),
        tTok("/"),
        tTok("<"),
        tTok(">"),
        tTok(";"),
        tTok("if"),
        tInt("5"),
        tTok("<"),
        tInt("10"),
        tTok("{"),
        tTok("return"),
        tTok("true"),
        tTok(";"),
        tTok("}"),
        tTok("else"),
        tTok("{"),
        tTok("return"),
        tTok("false"),
        tTok(";"),
        tTok("}"),
        tInt("5"),
        tTok("=="),
        tInt("5"),
        tTok(";"),
        tInt("5"),
        tTok("!="),
        tInt("10"),
        tTok(";"),
    }

    l := New(input)

    for i, tt := range tests {
        tok := l.NextToken()
        if tok.Type != tt.expType {
            t.Fatalf("token: %d - token type wrong. Expected %q, got %q",
                i + 1, tt.expType, tok.Type)
        }
        if tok.Literal != tt.expLiteral {
            t.Fatalf("token: %d - token literal wrong. Expected %q, got %q",
                i, tt.expLiteral, tok.Literal)
        }
    }
}

func tUser(l string) tToken {
    return tToken{expType: token.Ident, expLiteral: l}
}

func tInt(l string) tToken {
    return tToken{expType: token.Int, expLiteral: l }
}

func tTok(l string) (t tToken) {
    t.expLiteral = l;
    switch l {
        case "\x00": t.expType = token.EOF
        case ",": t.expType = token.Comma
        case ";": t.expType = token.Semicolon
        case "(": t.expType = token.LParen
        case ")": t.expType = token.RParen
        case "{": t.expType = token.LBrace
        case "}": t.expType = token.RBrace
        case "=": t.expType = token.Assign
        case "+": t.expType = token.Plus
        case "-": t.expType = token.Minus
        case "!": t.expType = token.Bang
        case "*": t.expType = token.Asterisk
        case "/": t.expType = token.Slash
        case "<": t.expType = token.LT
        case ">": t.expType = token.GT
        case "==": t.expType = token.EQ
        case "!=": t.expType = token.NOTEQ
        case "fn": t.expType = token.Function
        case "let": t.expType = token.Let
        case "true": t.expType = token.True
        case "false": t.expType = token.False
        case "if": t.expType = token.If
        case "else": t.expType = token.Else
        case "return": t.expType = token.Return
    }
    return
}
