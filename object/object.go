package object

import (
    "fmt"
    "strings"

    "lemur/ast"
)


type Object interface {
    Type()    ObjectType
    String() string
}

type ObjectType string // this can be a numeric enum

const (
    BuiltinType  = "Builtin"
    FunctionType = "Function"
    ArrayType	 = "Array"
    StringType	 = "String"
    IntegerType  = "Integer"
    BooleanType  = "Boolean"
    NullType     = "Null"
    ReturnType   = "Return"
    ErrorType    = "Error"
)

type Builtin func(args ...Object) Object

func (b Builtin) Type() ObjectType { return BuiltinType }
func (b Builtin) String() string { return "builtin function" }

type Function struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    OuterEnv   *Environment
}
var _ Object = (*Function)(nil)

func (f *Function) Type() ObjectType { return FunctionType }
func (f *Function) String() string {
    var out strings.Builder

    params := []string{}
    for _, p := range f.Parameters {
	    params = append(params, p.String())
    }

    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(")")
    out.WriteString(f.Body.String())

    return out.String()
}

type Array struct {
    Elements []Object
}
var _ Object = (*Array)(nil)

func (a *Array) Type() ObjectType { return ArrayType }
func (a *Array) String() string {
    var out strings.Builder

    elems := []string{}
    for _, el := range a.Elements {
	elems = append(elems, el.String())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elems, ", "))
    out.WriteString("]")

    return out.String()
}

type String struct {
    Value string
}
var _ Object = (*String)(nil)

func (s *String) Type() ObjectType { return StringType }
func (s *String) String() string { return s.Value }

type Integer struct {
    Value int64
}
var _ Object = (*Integer)(nil)

func (i *Integer) Type() ObjectType { return IntegerType }
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
    Value bool
}
var _ Object = (*Boolean)(nil)

func (b *Boolean) Type() ObjectType { return BooleanType }
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

type Null struct { // replace with sum type (option)?
    Value bool
}
var _ Object = (*Null)(nil)

func (b *Null) Type() ObjectType { return NullType }
func (b *Null) String() string { return "null" }

type Return struct {
    Value Object
}
var _ Object = (*Return)(nil)

func (r *Return) Type() ObjectType { return ReturnType }
func (r *Return) String() string { return r.Value.String() }

type Error struct { // attach token info to this
    Message string
}
var _ Object = (*Error)(nil)

func (e *Error) Type() ObjectType { return ErrorType }
func (e *Error) String() string { return "Error: " + e.Message }
