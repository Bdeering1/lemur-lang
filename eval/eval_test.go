package eval

import (
    "testing"

    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
)

func TestErrorCases(t *testing.T) {
    tests := []struct{
        input    string
        expected string
    }{
        {"!1", UnknownOperatorError + ": " + "!Integer"},
        {"!1; 2", UnknownOperatorError + ": " + "!Integer"},
        {"-true", UnknownOperatorError + ": " + "-Boolean"},
        {"-true; 2", UnknownOperatorError + ": " + "-Boolean"},
        {"true + true", UnknownOperatorError + ": " + "Boolean + Boolean"},
        {"true + true; 2", UnknownOperatorError + ": " + "Boolean + Boolean"},
        {"1 + true", TypeMismatchError + ": " + "Integer + Boolean"},
        {"true + 1", TypeMismatchError + ": " + "Boolean + Integer"},
        {"1 + true; 2", TypeMismatchError + ": " + "Integer + Boolean"},
        {"if 1 + 1 { 2 }", InvalidConditionError + ": " + "(1 + 1)"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Error](t, i, obj)
        assert(t, i, res.Message, tst.expected)
    }
}

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

        ret := assertCast[*object.Return](t, i, obj)
        n := assertCast[*object.Integer](t, i, ret.Value)
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
        res := assertCast[*object.Integer](t, i, obj)
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

        res := assertCast[*object.Integer](t, i, obj)
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

        res := assertCast[*object.Boolean](t, i, obj)
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
        t.Errorf("test %d: incorrect object value, expected %T: %v (got %T: %v)",
            testIdx + 1,
            expected, expected,
            val, val)
    }
}

func assertCast[T object.Object](t *testing.T, testIdx int, obj object.Object) T {
    o, ok := obj.(T)
    if !ok {
        t.Fatalf("test %d: object is not an %T (got %T)", testIdx + 1, *new(T), obj)
    }

    return o
}
