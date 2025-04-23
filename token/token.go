package token

type TokenType string // could make this int or byte

type Token struct {
    Type TokenType
    Literal string
}

const (
    Illegal = "Illegal"
    EOF     = "EOF"

    // Identifiers & Literals
    Ident = "Identifier"
    Int   = "Int"

    // Delimiters
    Comma     = ","
    Semicolon = ";"
    LParen    = "("
    RParen    = ")"
    LBrace    = "{"
    RBrace    = "}"

    // Operators
    Assign   = "="
    Plus     = "+"
    Minus    = "-"
    Bang     = "!"
    Asterisk = "*"
    Slash    = "/"

    LT       = "<"
    GT       = ">"
    EQ       = "=="
    NOTEQ    = "!="

    // Keywords
    Function = "Function"
    Let      = "Let"
    True     = "True"
    False    = "False"
    If       = "If"
    Else     = "Else"
    Return   = "Return"
)

var Operators = map[string]TokenType{
    "=":  Assign,
    "+":  Plus,
    "-":  Minus,
    "!":  Bang,
    "*":  Asterisk,
    "/":  Slash,
    "<":  LT,
    ">":  GT,
    "==": EQ,
    "!=": NOTEQ,
}

var Keywords = map[string]TokenType{
    "fn": Function,
    "let": Let,
    "true": True,
    "false": False,
    "if": If,
    "else": Else,
    "return": Return,
}

func OperatorType(op string) TokenType {
    if ot, ok := Operators[op]; ok {
        return ot
    }
    return Illegal
}

func IdentType(ident string) TokenType {
    if tt, ok := Keywords[ident]; ok {
        return tt
    }
    return Ident
}

func New(ttype TokenType, literal string) Token {
    return Token{Type: ttype, Literal: literal}
}
