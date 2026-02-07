package ast

import (
    "bytes"
    "fmt"
    "reflect"
    "strings"

    "lemur/token"
)

type Node interface {
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


type Program []Statement
var _ Node = (Program)(nil)

func (p Program) String() string {
    var out bytes.Buffer

    for _, s := range p {
        out.WriteString(s.String())
    }
    return out.String()
}
func (p *Program) PrintAST() string {
    var b strings.Builder

    prettyPrint(&b, reflect.ValueOf(*p), 0)
    b.WriteString("\n")

    return b.String()
}

func prettyPrint(b *strings.Builder, val reflect.Value, indent int) {
    indentStr := strings.Repeat("  ", indent)

    k := val.Kind()
    if k == reflect.Ptr || k == reflect.Interface {
		if val.IsNil() { b.WriteString("<nil>"); return }

        prettyPrint(b, val.Elem(), indent)
        return
    }

    switch k {
    case reflect.Struct:
	    t := val.Type()

	    b.WriteString(fmt.Sprintf("%s {\n", t.Name()))
	    for i := range val.NumField() {
		    b.WriteString(fmt.Sprintf("%s  %s: ", indentStr, t.Field(i).Name))
		    prettyPrint(b, val.Field(i), indent + 1)
		    b.WriteString("\n")
	    }
	    b.WriteString(fmt.Sprintf("%s}", indentStr))

    case reflect.Slice, reflect.Array:
	    if val.Len() == 0 { b.WriteString("[]"); return }

	    b.WriteString("[\n")
	    for i := range val.Len() {
	b.WriteString(fmt.Sprintf("%s  ", indentStr))
		    prettyPrint(b, val.Index(i), indent + 1)
		    b.WriteString(",\n")
	    }
	    b.WriteString(fmt.Sprintf("%s]", indentStr))

    case reflect.String:
	    b.WriteString(fmt.Sprintf("%q", val.String()))
    default:
	    b.WriteString(fmt.Sprintf("%+v", val.Interface()))
    }
}

type BlockStatement struct {
    Token        token.Token
    Statements []Statement
}
var _ Statement = (*BlockStatement)(nil)

func (bs *BlockStatement) _stmtNode(){}
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
func (ls *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString(ls.Token.Literal + " ")
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
func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.Token.Literal + " ")
    out.WriteString(rs.Value.String())
    out.WriteString(";")

    return out.String()
}

type ExpressionStatement struct {
    Token token.Token
    Value Expression
}
var _ Statement = (*ExpressionStatement)(nil);

func (es *ExpressionStatement) _stmtNode(){}
func (es *ExpressionStatement) String() string {
    var out bytes.Buffer

    out.WriteString(es.Value.String())
    out.WriteString(";")

    return out.String()
}

type Identifier struct {
    Token token.Token
    Value string
}
var _ Expression = (*Identifier)(nil)

func (i *Identifier) _exprNode(){}
func (i *Identifier) String() string { return i.Value }

type StringLiteral struct {
    Token token.Token
    Value string
}
var _ Expression = (*StringLiteral)(nil)

func (ll *StringLiteral) _exprNode(){}
func (sl *StringLiteral) String() string { return sl.Token.Literal }

type IntegerLiteral struct {
    Token token.Token
    Value int64
}
var _ Expression = (*IntegerLiteral)(nil)

func (il *IntegerLiteral) _exprNode(){}
func (il *IntegerLiteral) String() string { return il.Token.Literal }

type BooleanLiteral struct {
    Token token.Token
    Value bool
}
var _ Expression = (*BooleanLiteral)(nil)

func (b *BooleanLiteral) _exprNode(){}
func (b *BooleanLiteral) String() string { return b.Token.Literal }

type PrefixExpression struct {
    Token    token.Token
    Operator string
    Right    Expression
}
var _ Expression = (*PrefixExpression)(nil)

func (pe *PrefixExpression) _exprNode(){}
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
    Token	 token.Token
    Condition    Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}
var _ Expression = (*ConditionalExpression)(nil)

func (ce *ConditionalExpression) _exprNode(){}
func (ce *ConditionalExpression) String() string {
    var out bytes.Buffer

    out.WriteString("if ")
    out.WriteString(ce.Condition.String())
    out.WriteString(" ")
    out.WriteString(ce.Consequence.String())

    if (ce.Alternative != nil) {
        out.WriteString("else")
        out.WriteString(ce.Alternative.String())
    }

    return out.String()
}

type FunctionLiteral struct {
    Token	 token.Token
    Parameters []*Identifier
    Body	 *BlockStatement
}
var _ Expression = (*FunctionLiteral)(nil)

func (fl *FunctionLiteral) _exprNode(){}
func (fl *FunctionLiteral) String() string {
    var out bytes.Buffer

    params := []string{}
    for _, p := range fl.Parameters {
	params = append(params, p.String())
    }

    out.WriteString(fl.Token.Literal)
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(")")
    out.WriteString(fl.Body.String())

    return out.String()
}

type CallExpression struct {
    Token token.Token
    Function Expression // identifier (or function literal?)
    Arguments []Expression
}

func (ce *CallExpression) _exprNode(){}
func (ce *CallExpression) String() string {
    var out bytes.Buffer

    args := []string{}
    for _, a := range ce.Arguments {
	args = append(args, a.String())
    }

    out.WriteString(ce.Function.String())
    out.WriteString("(")
    out.WriteString(strings.Join(args, ", "))
    out.WriteString(")")

    return out.String()
}
