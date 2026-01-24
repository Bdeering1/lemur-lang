package eval

import (
    "testing"

    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input    string
        expected int64
    }{
        {"0", 0},
        {"15", 15},
    }

    for _, tst := range tests {
        e := runNewEval(tst.input)

        res := assertCast[*object.Integer](t, e)
        assert(t, res.Value, tst.expected)
    }
}

func runNewEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    return Eval(program)
}

func assert(t *testing.T, val any, expected any) {
    if val != expected {
        t.Errorf("incorrect object value, expected %v (got %v)", expected, val)
    }
}

func assertCast[T object.Object](t *testing.T, obj object.Object) T {
    o, ok := obj.(T)
    if !ok {
        t.Fatalf("object is not an %T (got %T)", *new(T), obj)
    }

    return o
}
