package ast

import (
    "bytes"

    "lemur/token"
)

type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface {
    Node
    _stmtNode()
}

type Expression interface {
    Node
    _exprNode()
}

type Program struct {
    Statements []Statement
}
var _ Node = (*Program)(nil)

func (p *Program) TokenLiteral() string {
    if len(p.Statements) == 0 { return "" }
    return p.Statements[0].TokenLiteral()
}
func (p *Program) String() string {
    var out bytes.Buffer

    for _, s := range p.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

type LetStatement struct {
    Token token.Token
    Name *Identifier // isn't this the same as the identifier token?
    Value Expression
}
var _ Statement = (*LetStatement)(nil)

func (ls *LetStatement) _stmtNode(){}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Name.String())
    out.WriteString(" = ")

    if ls.Value != nil { // temp. nil check
        out.WriteString(ls.Value.String())
    }
    out.WriteString(";")

    return out.String()
}

type ReturnStatement struct {
    Token token.Token
    Value Expression
}
var _ Statement = (*ReturnStatement)(nil)

func (rs *ReturnStatement) _stmtNode(){}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")

    if rs.Value != nil { // temp. nil check
        out.WriteString(rs.Value.String())
    }
    out.WriteString(";")

    return out.String()
}

type ExpressionStatement struct {
    Token token.Token
    Value Expression
}
var _ Statement = (*ExpressionStatement)(nil);

func (es *ExpressionStatement) _stmtNode(){}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
    var out bytes.Buffer

    out.WriteString(es.TokenLiteral() + " ")

    if es.Value != nil { // temp. nil check
        out.WriteString(es.Value.String())
    }
    out.WriteString(";")

    return out.String()
}

type Identifier struct {
    Token token.Token
    Value string
}
var _ Expression = (*Identifier)(nil)

func (i *Identifier) _exprNode(){}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }

type IntegerLiteral struct {
    Token token.Token
    Value int64
}
var _ Expression = (*IntegerLiteral)(nil)

func (il *IntegerLiteral) _exprNode(){}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.Token.Literal }
