package object

import "fmt"

type ObjectType string // this can be a numeric enum

type Object interface {
    Type()    ObjectType
    Inspect() string // could just be String() method
}

const (
    IntegerType = "Integer"
    BooleanType = "Boolean"
    NullType    = "Null"
)

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { return IntegerType }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanType }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct { // replace with sum type (option)?
    Value bool
}

func (b *Null) Type() ObjectType { return NullType }
func (b *Null) Inspect() string { return "null" }
