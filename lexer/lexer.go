package lexer

import "monkey/token"

type Lexer struct {
    input   string
    pos     int
    readPos int // is this needed?
    ch      byte
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    switch l.ch {
        case '=': tok.Type = token.Assign
        case '+': tok.Type = token.Plus
        case ',': tok.Type = token.Comma
        case ';': tok.Type = token.Semicolon
        case '(': tok.Type = token.LParen
        case ')': tok.Type = token.RParen
        case '{': tok.Type = token.LBrace
        case '}': tok.Type = token.RBrace
        case 0:   tok.Type = token.EOF
    }
    if l.ch == 0 { tok.Literal = "" } else { tok.Literal = string(l.ch) }

    l.readChar()
    return tok
}

func (l *Lexer) readChar() {
    if l.readPos >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPos]
    }
    l.pos = l.readPos
    l.readPos++
}
