package eval

import (
    "testing"

    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
)

func TestReturnStatement(t *testing.T) {
    tests := []struct{
        input    string
        expected int64
    }{
        {"return 10", 10},
        {"return 10; 9", 10},
        {"return 2 * 5; 9", 10},
        {"8; return 2 * 5; 9", 10},
        {"{ return 10; }", 10},
        {"{ return 10; 9 }", 10},
        {"{{ return 10; 9 } 8 }", 10},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        ret := assertCast[*object.Return](t, obj)
        n := assertCast[*object.Integer](t, ret.Value)
        assert(t, i, n.Value, tst.expected)
    }
}

func TestConditionalExpression(t *testing.T) {
    tests := []struct{
        input    string
        expected any
    }{
        {"if true { 10 }", 10},
        {"if false { 10 }", nil},
        {"if 1 < 2 { 10 }", 10},
        {"if 1 > 2 { 10 }", nil},
        {"if true { 10 } else { 20 }", 10},
        {"if false { 10 } else { 20 }", 20},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)
        expd, ok := tst.expected.(int)

        if !ok {
            assert(t, i, obj, Null)
            continue
        }
        res := assertCast[*object.Integer](t, obj)
        assert(t, i, res.Value, int64(expd))
    }
}

func TestIntegerExpression(t *testing.T) {
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

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Integer](t, obj)
        assert(t, i, res.Value, tst.expected)
    }
}

func TestBooleanExpression(t *testing.T) {
    tests := []struct{
        input    string
        expected bool
    }{
        {"true", true},
        {"false", false},
        {"!true", false},
        {"!false", true},
        {"!!false", false},
        {"!!true", true},
        {"true == true", true},
        {"false == false", true},
        {"true == false", false},
        {"true != false", true},
        {"false != true", true},
        {"1 < 2", true},
        {"1 > 2", false},
        {"1 < 1", false},
        {"1 > 1", false},
        {"1 == 1", true},
        {"1 != 1", false},
        {"1 == 2", false},
        {"1 != 2", true},
        {"(1 < 2) == true", true},
        {"(1 < 2) == false", false},
        {"(1 > 2) == true", false},
        {"(1 > 2) == false", true},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Boolean](t, obj)
        assert(t, i, res.Value, tst.expected)
    }
}

func runNewEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    return Eval(program)
}

func assert(t *testing.T, testIdx int, val any, expected any) {
    if val != expected {
        t.Errorf("incorrect object value for test %d, expected %T: %v (got %T: %v)",
            testIdx + 1,
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
