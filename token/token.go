package token

type TokenType string // could make this int or byte

type Token struct {
    Type TokenType
    Literal string
}

const (
    Illegal = "Illegal"
    EOF = "EOF"

    // Identifiers + literals
    Ident = "Identifier"
    Int = "Int"

    // Operators
    Assign = "="
    Plus = "+"

    // Delimiters
    Comma = ","
    Semicolon = ";"
    LParen = "("
    RParen = ")"
    LBrace = "{"
    RBrace = "}"

    // Keywords
    Function = "Function"
    Let = "Let"
)
