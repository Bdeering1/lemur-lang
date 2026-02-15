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

func TestBuiltinFunction(t *testing.T) { // these should use builtin constants
    tests := []struct{
        input    string
        expected any
    }{
        {"len([])", 0},
        {"len([1, 2, 3])", 3},
        {`len("")`, 0},
        {`len("four")`, 4},
        {`len("1")`, 1},
        {`len(1)`, ArgumentTypesError + ": len(Integer)"},
        {`len(true)`, ArgumentTypesError + ": len(Boolean)"},
        {`len([], [])`, ArgumentMistmatchError + ": len"},
        {"first([])", nil},
        {"first([1, 2, 3])", 1},
        {"first(1)", ArgumentTypesError + ": first(Integer)"},
        {"first(true)", ArgumentTypesError + ": first(Boolean)"},
        {`first([], [])`, ArgumentMistmatchError + ": first"},
        {"last([])", nil},
        {"last([1, 2, 3])", 3},
        {"last(1)", ArgumentTypesError + ": last(Integer)"},
        {"last(true)", ArgumentTypesError + ": last(Boolean)"},
        {`last([], [])`, ArgumentMistmatchError + ": last"},
        {"head([])", []int{}},
        {"head([1, 2, 3])", []int{1, 2}},
        {"head(1)", ArgumentTypesError + ": head(Integer)"},
        {"head(true)", ArgumentTypesError + ": head(Boolean)"},
        {`head([], [])`, ArgumentMistmatchError + ": head"},
        {"tail([])", []int{}},
        {"tail([1, 2, 3])", []int{2, 3}},
        {"tail(1)", ArgumentTypesError + ": tail(Integer)"},
        {"tail(true)", ArgumentTypesError + ": tail(Boolean)"},
        {`tail([], [])`, ArgumentMistmatchError + ": tail"},
        {"push([], 1)", []int{1}},
        {"push([1, 2], 3)", []int{1, 2, 3}},
        {"push([1, 2], true)", TypeMismatchError + ": push(Array[Integer], Boolean)"},
        {"push([true, false], 1)", TypeMismatchError + ": push(Array[Boolean], Integer)"},
        {"push(1, 1)", ArgumentTypesError + ": push(Integer, Integer)"},
        {"push(true, true)", ArgumentTypesError + ": push(Boolean, Boolean)"},
        {`push([])`, ArgumentMistmatchError + ": push"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[*object.Integer](t, i, obj)
            assert(t, i, res.Value, int64(expd))
        case []int:
            arr := assertCast[*object.Array](t, i, obj)
            for idx, el := range arr.Elements {
                res := el.(*object.Integer)
                assert(t, i, res.Value, int64(expd[idx]))
            }
        case string:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, expd)
        case nil:
            assert(t, i, obj, Null)
        }
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

        assertMsg(t, i, len(f.Parameters), tst.expdParams, "wrong nunmber of parameters in function object")
        assertMsg(t, i, f.Body.String(), tst.expdBody, "incorrect function body")
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

func TestArrayLiteral(t *testing.T) {
    input := "[1, 2 * 3, fn(){ 5 * 6 }()]"

    obj := runNewEval(input)
    arr := assertCast[*object.Array](t, 0, obj)

    first := assertCast[*object.Integer](t, 0, arr.Elements[0])
    assert(t, 0, first.Value, int64(1))
    second := assertCast[*object.Integer](t, 0, arr.Elements[1])
    assert(t, 0, second.Value, int64(6))
    third:= assertCast[*object.Integer](t, 0, arr.Elements[2])
    assert(t, 0, third.Value, int64(30))
}

func TestIndexExpression(t *testing.T) {
    tests := []struct{
        input    string
        expected any
    }{
        {"[1, 2][0]", 1},
        {"[1, 2][0 + 1]", 2},
        {"let arr = [1, 2, 3]; arr[2]", 3},
        {`"hello"[0]`, "h"},
        {`"world"[1]`, "o"},
        {`let s = "asdf"; s[2]`, "d"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[*object.Integer](t, i, obj)
            assert(t, i, res.Value, int64(expd))
        case string:
            res := assertCast[*object.String](t, i, obj)
            assert(t, i, res.Value, expd)
        }
    }
}

func TestStringLiteral(t *testing.T) {
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
        {"true && true", true},
        {"true && false", false},
        {"false && true", false},
        {"false && false", false},
        {"true || true", true},
        {"true || false", true},
        {"false || true", true},
        {"false || false", false},
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
        {"!1", UnknownOperatorError + ": !Integer"},
        {"!1; 2", UnknownOperatorError + ": !Integer"},
        {"-true", UnknownOperatorError + ": -Boolean"},
        {"-true; 2", UnknownOperatorError + ": -Boolean"},
        {"true + true", UnknownOperatorError + ": Boolean + Boolean"},
        {"true + true; 2", UnknownOperatorError + ": Boolean + Boolean"},
        {`"foo" - "bar"`, UnknownOperatorError + ": String - String"},
        {"1 && 0", UnknownOperatorError + ": Integer && Integer"},
        {`"a" && "b"`, UnknownOperatorError + ": String && String"},
        {"1 || 0", UnknownOperatorError + ": Integer || Integer"},
        {`"a" || "b"`, UnknownOperatorError + ": String || String"},

        {"1 + true", TypeMismatchError + ": Integer + Boolean"},
        {"true + 1", TypeMismatchError + ": Boolean + Integer"},
        {"!(true + 1)", TypeMismatchError + ": Boolean + Integer"},
        {"(true + 1) * (5 + 5)", TypeMismatchError + ": Boolean + Integer"},
        {"if true + 1 { 2 }", TypeMismatchError + ": Boolean + Integer"},
        {"return true + 1", TypeMismatchError + ": Boolean + Integer"},
        {"1 + true; 2", TypeMismatchError + ": Integer + Boolean"},
        {"true && 1", TypeMismatchError + ": Boolean && Integer"},
        {"0 && false", TypeMismatchError + ": Integer && Boolean"},
        {"true || 1", TypeMismatchError + ": Boolean || Integer"},
        {"0 || false", TypeMismatchError + ": Integer || Boolean"},

        {"if 1 + 1 { 2 }", InvalidConditionError + ": (1 + 1)"},

        {"x", IdentifierNotFoundError + ": x"},
        {"!x", IdentifierNotFoundError + ": x"},
        {"if x { y }", IdentifierNotFoundError + ": x"},
        {"return x", IdentifierNotFoundError + ": x"},

        {"[1, 2][-1]", IndexOutOfBoundsError + ": -1"},
        {"[1, 2][2]", IndexOutOfBoundsError + ": 2"},
        {`"hello"[-1]`, IndexOutOfBoundsError + ": -1"},
        {`"world"[5]`, IndexOutOfBoundsError + ": 5"},

        {"[1, 2][true]", InvalidIndexExpressionError + ": cannot index Array with Boolean"},
        {`[1, 2]["asdf"]`, InvalidIndexExpressionError + ": cannot index Array with String"},
        {`""[true]`, InvalidIndexExpressionError + ": cannot index String with Boolean"},
        {`""["asdf"]`, InvalidIndexExpressionError + ": cannot index String with String"},
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

func assert(t *testing.T, testIdx int, val, expected any) {
    if val != expected {
        t.Errorf("test %d: incorrect object value, expected %T: %v (got %T: %v)",
            testIdx + 1,
            expected, expected,
            val, val)
    }
}

func assertMsg(t *testing.T, testIdx int, val, expected any, msg string) {
    if val != expected {
        t.Fatalf("test %d: %s, expected %T: %v (got %T: %v)",
            testIdx + 1,
            msg,
            expected, expected,
            val, val)
    }
}

func assertCast[T object.Object](t *testing.T, testIdx int, obj object.Object) T {
    o, ok := obj.(T)
    if !ok {
        if isError(obj) { t.Errorf("%s", obj.String()) }
        t.Fatalf("test %d: object is not an %T (got %T)", testIdx + 1, *new(T), obj)
    }

    return o
}
