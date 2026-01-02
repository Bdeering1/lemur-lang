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

type BlockStatement struct {
    Token        token.Token
    Statements []Statement
}
var _ Statement = (*BlockStatement)(nil)

func (bs *BlockStatement) _stmtNode(){}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    out.WriteString("{")
    for _, s := range bs.Statements {
        out.WriteString(s.String())
    }
    out.WriteString("}")

    return out.String()
}

type LetStatement struct {
    Token token.Token
    Name *Identifier
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

type BooleanLiteral struct {
    Token token.Token
    Value bool
}
var _ Expression = (*BooleanLiteral)(nil)

func (b *BooleanLiteral) _exprNode(){}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string { return b.Token.Literal }

type PrefixExpression struct {
    Token    token.Token
    Operator string
    Right    Expression
}
var _ Expression = (*PrefixExpression)(nil)

func (pe *PrefixExpression) _exprNode(){}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

type InfixExpression struct {
    Token    token.Token
    Operator string
    Left     Expression
    Right    Expression
}
var _ Expression = (*InfixExpression)(nil)

func (ie *InfixExpression) _exprNode(){}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")

    return out.String()
}

type ConditionalExpression struct {
    Token token.Token
    Condition    Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}
var _ Expression = (*ConditionalExpression)(nil)

func (ce *ConditionalExpression) _exprNode(){}
func (ce *ConditionalExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ConditionalExpression) String() string {
    var out bytes.Buffer

    out.WriteString("if")
    out.WriteString(ce.Condition.String())
    out.WriteString(" ")
    out.WriteString(ce.Consequence.String())

    if (ce.Alternative != nil) {
        out.WriteString("else")
        out.WriteString(ce.Alternative.String())
    }

    return out.String()
}
