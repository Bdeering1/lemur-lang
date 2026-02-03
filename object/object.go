package object

import (
    "bytes"
	"fmt"
    "strings"

	"lemur/ast"
)

type ObjectType string // this can be a numeric enum

type Object interface {
    Type()    ObjectType
    String() string
}

const (
    FunctionType = "Function"
    IntegerType  = "Integer"
    BooleanType  = "Boolean"
    NullType     = "Null"
    ReturnType   = "Return"
    ErrorType    = "Error"
)

type Function struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    OuterEnv   *Environment
}

func (f *Function) Type() ObjectType { return IntegerType }
func (f *Function) String() string {
    var out bytes.Buffer

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


type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { return IntegerType }
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanType }
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

type Null struct { // replace with sum type (option)?
    Value bool
}

func (b *Null) Type() ObjectType { return NullType }
func (b *Null) String() string { return "null" }

type Return struct {
    Value Object
}

func (r *Return) Type() ObjectType { return ReturnType }
func (r *Return) String() string { return r.Value.String() }

type Error struct { // attach token info to this
    Message string
}

func (e *Error) Type() ObjectType { return ErrorType }
func (e *Error) String() string { return "Error: " + e.Message }
