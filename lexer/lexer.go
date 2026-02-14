package lexer

import (
    "fmt"

    "lemur/token"
)

type Lexer struct {
    input   string
    pos     int
    nextPos int
    ch      byte
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) NextToken() (tok token.Token) {
    l.skipWhitespace()
    tok.Literal = string(l.ch)

    switch l.ch {
    case '\x00': tok.Type = token.EOF
    case ',': tok.Type = token.Comma
    case ';': tok.Type = token.Semicolon
    case '(': tok.Type = token.LParen
    case ')': tok.Type = token.RParen
    case '{': tok.Type = token.LBrace
    case '}': tok.Type = token.RBrace
    case '[': tok.Type = token.LBracket
    case ']': tok.Type = token.RBracket
    case '+': tok.Type = token.Plus
    case '-': tok.Type = token.Minus
    case '*': tok.Type = token.Asterisk
    case '/': tok.Type = token.Slash
    case '>': tok.Type = token.GT
    case '<': tok.Type = token.LT
    case '=', '!', '&', '|':
        tok.Literal = l.readOperator()
        tok.Type = token.OperatorType(tok.Literal)
    case '"':
        tok.Type = token.String
        tok.Literal = l.readString()
    default:
        if isAlpha(l.ch) {
            tok.Literal = l.readIdent()
            tok.Type = token.IdentType(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            l.readNumber(&tok)
            return tok
        }
        fmt.Printf("Illegal char: %d", l.ch)
        tok.Type = token.Illegal
    }

    l.readChar()
    return tok
}

func (l *Lexer) readString() string {
    startPos := l.pos + 1
    for {
        l.readChar()
        if l.ch == '"' || l.ch == '\x00' { break }
    }

    return l.input[startPos : l.pos]
}

func (l *Lexer) readOperator() string { // make this match readNumber
    cur := string(l.ch)
    literal := string(cur) + string(l.nextChar())
    if isOperator(literal) {
        l.readChar()
        return literal
    }

    if isOperator(cur) { return cur }
    return token.Illegal
}

func isOperator(op string) bool {
    _, ok := token.Operators[op]; return ok
}

func (l *Lexer) readIdent() string {
    pos := l.pos
    for isAlpha(l.ch) || isDigit(l.ch) {
        l.readChar()
    }
    return l.input[pos:l.pos]
}

func (l *Lexer) readNumber(tok *token.Token) {
    pos := l.pos

    valid := true
    for isDigit(l.ch) || isAlpha(l.ch) {
        if !isDigit(l.ch) { valid = false }
        l.readChar()
    }
    tok.Literal = l.input[pos:l.pos]

    if !valid {
        tok.Type = token.Illegal
        return
    }
    tok.Type = token.Int
}


func (l *Lexer) readChar() {
    l.pos = l.nextPos
    l.nextPos++

    if l.pos >= len(l.input) {
        l.ch = '\x00'
        return
    }
    l.ch = l.input[l.pos]
}

func (l *Lexer) nextChar() byte {
    if l.nextPos >= len(l.input) {
        return '\x00'
    }
    return l.input[l.nextPos]
}

func (l *Lexer) skipWhitespace() {
    for l.charIsWhiteSpace() {
        l.readChar()
    }
}

func (l *Lexer) charIsWhiteSpace() bool {
    return l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r'
}

func isAlpha(ch byte) bool {
    return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}
