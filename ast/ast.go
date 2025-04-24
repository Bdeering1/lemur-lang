package ast

import "lemur/token"

type Node interface {
    TokenLiteral() string // can this be a dummy method?
}

type Statement interface {
    Node
    stmtNode() // dummy method
}

type Expression interface {
    Node
    exprNode() // dummy method
}

type Program struct {
    Statements []Statement
}
var _ Node = (*Program)(nil)

func (p *Program) TokenLiteral() string {
    if len(p.Statements) == 0 { return "" }
    return p.Statements[0].TokenLiteral() // why?
}

type Identifier struct {
    Token token.Token
    Value string
}
var _ Expression = (*Identifier)(nil)

func (i *Identifier) exprNode(){}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type LetStatement struct {
    Token token.Token // token.Let (maybe only needed if the token contained line/col info)
    Name *Identifier // can this just be the token itself?
    Value Expression
}
var _ Statement = (*LetStatement)(nil)

func (ls *LetStatement) stmtNode(){}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
