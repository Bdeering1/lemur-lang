package ast

import "lemur/token"

type Node interface {
    TokenLiteral() string // can this be a dummy method?
}

type Statement interface {
    Node
    _stmtNode() // dummy method
}

type Expression interface {
    Node
    _exprNode() // dummy method
}

type Program struct {
    Statements []Statement
}
var _ Node = (*Program)(nil)

func (p *Program) TokenLiteral() string {
    if len(p.Statements) == 0 { return "" }
    return p.Statements[0].TokenLiteral() // why?
}

type LetStatement struct {
    Token token.Token // token.Let (maybe only needed if the token contained line/col info)
    Name *Identifier // isn't this the same as the identifier token?
    Value Expression
}
var _ Statement = (*LetStatement)(nil)

func (ls *LetStatement) _stmtNode(){}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
    Token token.Token
    Value Expression
}
var _ Statement = (*ReturnStatement)(nil)

func (rs *ReturnStatement) _stmtNode(){}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type Identifier struct {
    Token token.Token
    Value string
}
var _ Expression = (*Identifier)(nil)

func (i *Identifier) _exprNode(){}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
