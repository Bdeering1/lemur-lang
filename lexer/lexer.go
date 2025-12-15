package lexer

import ("lemur/token")

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
    case '=', '+', '-', '!', '*', '/', '<', '>': // checking individually would be more performant
        tok.Literal = l.readOperator()
        tok.Type = token.OperatorType(tok.Literal)
    default:
        if isAlpha(l.ch) {
            tok.Literal = l.readIdent()
            tok.Type = token.IdentType(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            tok.Literal = l.readNumber()
            tok.Type = token.Int
            return tok
        }
        tok.Type = token.Illegal
    }

    l.readChar()
    return tok
}

func (l *Lexer) readOperator() string {
    literal := string(l.ch) + string(l.peekChar())
    if isOperator(literal) {
        l.readChar()
        return literal
    }
    return string(l.ch)
}

func isOperator(op string) bool {
    _, ok := token.Operators[op]; return ok
}

func (l *Lexer) readIdent() string {
    pos := l.pos
    for isAlpha(l.ch) { // support numeric characters?
        l.readChar()
    }
    return l.input[pos:l.pos]
}

func isAlpha(ch byte) bool {
    return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
    pos := l.pos
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[pos:l.pos]
}

func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
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

func (l *Lexer) peekChar() byte {
    if l.nextPos >= len(l.input) {
        return '\x00'
    }
    return l.input[l.nextPos]
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}
