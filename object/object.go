package object

import (
    "fmt"
    "strings"

    "lemur/ast"
)


type Object interface {
    Type()   ObjectType
    String() string
}

type ObjectType string // this can be a numeric enum

const (
    TupleType	 = "Tuple"
    BuiltinType  = "Builtin"
    FunctionType = "Function"
    HashMapType	 = "HashMap"
    ArrayType	 = "Array"
    StringType	 = "String"
    IntegerType  = "Integer"
    BooleanType  = "Boolean"
    ReturnType   = "Return"
    ErrorType    = "Error"
    NullType     = ""
)

type ValueObject interface { _valueObject() }


func Default(t ObjectType) Object {
    switch t {
    case ArrayType:
	return &Array{Elements: []Object{}}
    case StringType:
	return String("")
    case IntegerType:
	return Integer(0)
    case BooleanType:
	return Boolean(false)
    default:
	return NullObject(false)
    }
}


type Tuple []Object
var _ Object = (Tuple)(nil)

func (t Tuple) Type() ObjectType { return TupleType }
func (t Tuple) String() string {
    var out strings.Builder

    elems := []string{}
    for _, el := range t {
	elems = append(elems, el.String())
    }

    out.WriteString("(")
    out.WriteString(strings.Join(elems, ", "))
    out.WriteString(")")

    return out.String()
}

type Builtin func(args []Object, env *Environment) Object
var _ Object = (Builtin)(nil)
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

type HashMap struct {
    Pairs     map[Object]Object
    KeyType   ObjectType
    ValueType ObjectType
}
var _ Object = (*HashMap)(nil)

func (hm *HashMap) Type() ObjectType{ return HashMapType }
func (hm *HashMap) String() string {
    var out strings.Builder

    pairs := []string{}
    for k, v := range hm.Pairs {
	pairs = append(pairs, fmt.Sprintf("%s: %s", k, v))
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString("}")

    return out.String()
}

type Array struct {
    Elements    []Object
    ElementType ObjectType
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

type String string
var _ Object = (String)("")
var _ ValueObject = (String)("")
func (s String) Type() ObjectType { return StringType }
func (s String) String() string { return string(s) }
func (s String) _valueObject(){}

type Integer int
var _ Object = (Integer)(0)
var _ ValueObject = (Integer)(0)
func (i Integer) Type() ObjectType { return IntegerType }
func (i Integer) String() string { return fmt.Sprintf("%d", i) }
func (i Integer) _valueObject(){}

type Boolean bool
var _ Object = (Boolean)(false)
var _ ValueObject = (Boolean)(false)
func (b Boolean) Type() ObjectType { return BooleanType }
func (b Boolean) String() string { return fmt.Sprintf("%t", b) }
func (b Boolean) _valueObject(){}

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

type NullObject bool
var _ Object = (NullObject)(false)
func (b NullObject) Type() ObjectType { return NullType }
func (b NullObject) String() string { return "null" }

var Null = NullObject(false)
