package eval

import (
    "testing"

    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
)

func TestBangOperator(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"!true", false},
        {"!false", true},
        {"!!false", false},
        {"!!true", true},
    }

    for _, tst := range tests {
        e := runNewEval(tst.input)

        res := assertCast[*object.Boolean](t, e)
        assert(t, res.Value, tst.expected)
    }
}

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input    string
        expected int64
    }{
        {"0", 0},
        {"5", 5},
        {"10", 10},
        {"-0", 0},
        {"-5", -5},
        {"-10", -10},
        {"5 + 5 + 5", 15},
        {"20 - 5 - 5", 10},
        {"2 * 2 * 2", 8},
        {"20 / 2 / 2", 5},
        {"2 * (2 + 3)", 10},
        {"-7 + 7 + -7", -7},
        {"5 * 2 + 10", 20},
        {"10 + 5 * 2", 20},
    }

    for _, tst := range tests {
        e := runNewEval(tst.input)

        res := assertCast[*object.Integer](t, e)
        assert(t, res.Value, tst.expected)
    }
}

func TestEvalBooleanExpression(t *testing.T) {
    tests := []struct{
        input    string
        expected bool
    }{
        {"true", true},
        {"false", false},
    }

    for _, tst := range tests {
        e := runNewEval(tst.input)

        res := assertCast[*object.Boolean](t, e)
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
        t.Errorf("incorrect object value, expected %T: %v (got %T: %v)",
            expected, expected,
            val, val)
    }
}

func assertCast[T object.Object](t *testing.T, obj object.Object) T {
    o, ok := obj.(T)
    if !ok {
        t.Fatalf("object is not an %T (got %T)", *new(T), obj)
    }

    return o
}
