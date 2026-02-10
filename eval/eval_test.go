package eval

import (
    "testing"

    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
)

func TestLetStatement(t *testing.T) {
    tests := []struct{
        input    string
        expected int64
    }{
        {"let a = 5; a", 5},
        {"let a = 2 + 3; a", 5},
        {"let a = 5; let b = a; b", 5},
        {"let a = 2; let b = 3; a + b", 5},
    }


    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Integer](t, i, obj)
        assert(t, i, res.Value, tst.expected)
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

func TestFunctionExpression(t *testing.T) {
    tests := []struct{
        input      string
        expdParams int
        expdBody   string
    }{
        {"fn(x) { x }", 1, "{x;}"},
        {"fn(x, y) { x + 1 }", 2, "{(x + 1);}"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)
        f := assertCast[*object.Function](t, i, obj)

        if len(f.Parameters) != tst.expdParams {
            t.Fatalf("test %d: wrong number of parameters in function object, expected %d (got %d)",
                i,
                tst.expdParams,
                len(f.Parameters))
        }

        if f.Body.String() != tst.expdBody {
            t.Fatalf("test %d: incorrect function body, expected %s (got %s)",
                i,
                tst.expdBody,
                f.Body.String())
        }
    }
}

func TestCallExpression(t *testing.T) {
    tests := []struct{
        input    string
        expected int64
    }{
        {"let identity = fn(x) { x }; identity(5)", 5},
        {"let identity = fn(x) { return x }; identity(5)", 5},
        {"let double = fn(x) { x * 2 }; double(1)", 2},
        {"let add = fn(x, y) { x + y }; add(2, 3)", 5},
        {"let max = fn(x, y) { if x > y { x } else { y } }; max(1, 5)", 5},
        {"let fact = fn(n) { if n == 0 { 1 } else { n * fact(n-1) } }; fact(3)", 6},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Integer](t, i, obj)
        assert(t, i, res.Value, tst.expected)
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

func TestStringExpression(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {`"foo"`, "foo"},
        {`"Hello world!"`, "Hello world!"},
        {`"Hello" + " world!"`, "Hello world!"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.String](t, i, obj)
        assert(t, i, res.Value, tst.expected)
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
        {`"foo" == "foo"`, true},
        {`"foo" == "bar"`, false},
        {`"foo" != "foo"`, false},
        {`"foo" != "bar"`, true},
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
        {`"foo" - "bar"`, UnknownOperatorError + ": " + "String - String"},
        {"1 + true", TypeMismatchError + ": " + "Integer + Boolean"},
        {"true + 1", TypeMismatchError + ": " + "Boolean + Integer"},
        {"!(true + 1)", TypeMismatchError + ": " + "Boolean + Integer"},
        {"(true + 1) * (5 + 5)", TypeMismatchError + ": " + "Boolean + Integer"},
        {"if true + 1 { 2 }", TypeMismatchError + ": " + "Boolean + Integer"},
        {"return true + 1", TypeMismatchError + ": " + "Boolean + Integer"},
        {"1 + true; 2", TypeMismatchError + ": " + "Integer + Boolean"},
        {"if 1 + 1 { 2 }", InvalidConditionError + ": " + "(1 + 1)"},
        {"x", IdentifierNotFoundError + ": " + "x"},
        {"!x", IdentifierNotFoundError + ": " + "x"},
        {"if x { y }", IdentifierNotFoundError + ": " + "x"},
        {"return x", IdentifierNotFoundError + ": " + "x"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Error](t, i, obj)
        assert(t, i, res.Message, tst.expected)
    }
}

func runNewEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    env := object.CreateEnvironment()

    return Eval(program, env)
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
