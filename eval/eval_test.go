package eval

import (
    "testing"

    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
)

type Error string

func TestLetStatement(t *testing.T) {
    tests := []struct{
        input    string
        expected int
    }{
        {"let a = 5; a", 5},
        {"let a = 2 + 3; a", 5},
        {"let a = 5; let b = a; b", 5},
        {"let a = 2; let b = 3; a + b", 5},
    }


    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[object.Integer](t, i, obj)
        assert(t, i, res, object.Integer(tst.expected))
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct{
        input    string
        expected int
    }{
        {"return 10", 10},
        {"return 10; 9", 10},
        {"return 2 * 5; 9", 10},
        {"8; return 2 * 5; 9", 10},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        ret := assertCast[*object.Return](t, i, obj)
        n := assertCast[object.Integer](t, i, ret.Value)
        assert(t, i, n, object.Integer(tst.expected))
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
        {"len(1)", Error(ArgumentTypesError + ": len(Integer)")},
        {"len(true)", Error(ArgumentTypesError + ": len(Boolean)")},
        {"len([], [])", Error(ArgumentMistmatchError + ": len")},
        {"first([])", nil},
        {"first([1, 2, 3])", 1},
        {`first("")`, nil},
        {`first("asdf")`, "a"},
        {"first(1)", Error(ArgumentTypesError + ": first(Integer)")},
        {"first(true)", Error(ArgumentTypesError + ": first(Boolean)")},
        {"first([], [])", Error(ArgumentMistmatchError + ": first")},
        {"last([])", nil},
        {"last([1, 2, 3])", 3},
        {`last("")`, nil},
        {`last("asdf")`, "f"},
        {"last(1)", Error(ArgumentTypesError + ": last(Integer)")},
        {"last(true)", Error(ArgumentTypesError + ": last(Boolean)")},
        {`last([], [])`, Error(ArgumentMistmatchError + ": last")},
        {"head([])", []int{}},
        {"head([1, 2, 3])", []int{1, 2}},
        {`head("")`, ""},
        {`head("asdf")`, "asd"},
        {"head(1)", Error(ArgumentTypesError + ": head(Integer)")},
        {"head(true)", Error(ArgumentTypesError + ": head(Boolean)")},
        {"head([], [])", Error(ArgumentMistmatchError + ": head")},
        {"tail([])", []int{}},
        {"tail([1, 2, 3])", []int{2, 3}},
        {`tail("")`, ""},
        {`tail("asdf")`, "sdf"},
        {"tail(1)", Error(ArgumentTypesError + ": tail(Integer)")},
        {"tail(true)", Error(ArgumentTypesError + ": tail(Boolean)")},
        {"tail([], [])", Error(ArgumentMistmatchError + ": tail")},
        {"push([], 1)", []int{1}},
        {"push([1, 2], 3)", []int{1, 2, 3}},
        {"push([1, 2], true)", Error(TypeMismatchError + ": push(Array[Integer], Boolean)")},
        {"push([true, false], 1)", Error(TypeMismatchError + ": push(Array[Boolean], Integer)")},
        {"push(1, 1)", Error(ArgumentTypesError + ": push(Integer, Integer)")},
        {"push(true, true)", Error(ArgumentTypesError + ": push(Boolean, Boolean)")},
        {`push([])`, Error(ArgumentMistmatchError + ": push")},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[object.Integer](t, i, obj)
            assert(t, i, res, object.Integer(expd))
        case []int:
            arr := assertCast[*object.Array](t, i, obj)
            for idx, el := range arr.Elements {
                res := el.(object.Integer)
                assert(t, i, res, object.Integer(expd[idx]))
            }
        case Error:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, string(expd))
        case string:
            res := assertCast[object.String](t, i, obj)
            assert(t, i, res, object.String(expd))
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
        expected int
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

        res := assertCast[object.Integer](t, i, obj)
        assert(t, i, res, object.Integer(tst.expected))
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

        {"if 1 + 1 { 2 }", InvalidConditionError + ": (1 + 1)"},
        {`if "asdf" { 2 }`, InvalidConditionError + `: "asdf"`},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[object.Integer](t, i, obj)
            assert(t, i, res, object.Integer(expd))
        case Error:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, string(expd))
        case nil:
            assert(t, i, obj, Null)
        }
    }
}

func TestArrayLiteral(t *testing.T) {
    tests := []struct{
        input    string
        expected any
    }{
        {"[1, 2]", []int{1, 2}},
        {"[true, false]", []bool{true, false}},

        {"[1, true]", TypeMismatchError + ": Boolean in fixed-type Array of Integer"},
        {"[true, 1]", TypeMismatchError + ": Integer in fixed-type Array of Boolean"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case []int:
            arr := assertCast[*object.Array](t, i, obj)
            for idx, el := range arr.Elements {
                res := el.(object.Integer)
                assert(t, i, res, object.Integer(expd[idx]))
            }
        case []bool:
            arr := assertCast[*object.Array](t, i, obj)
            for idx, el := range arr.Elements {
                res := el.(object.Boolean)
                assert(t, i, res, object.Boolean(expd[idx]))
            }
        case Error:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, string(expd))
        }
    }
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

        {"[1, 2][-1]", Error(IndexOutOfBoundsError + ": -1")},
        {"[1, 2][2]", Error(IndexOutOfBoundsError + ": 2")},
        {`"hello"[-1]`, Error(IndexOutOfBoundsError + ": -1")},
        {`"world"[5]`, Error(IndexOutOfBoundsError + ": 5")},

        {"[1, 2][true]", Error(InvalidIndexExpressionError + ": cannot index Array with Boolean")},
        {`[1, 2]["asdf"]`, Error(InvalidIndexExpressionError + ": cannot index Array with String")},
        {`""[true]`, Error(InvalidIndexExpressionError + ": cannot index String with Boolean")},
        {`""["asdf"]`, Error(InvalidIndexExpressionError + ": cannot index String with String")},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[object.Integer](t, i, obj)
            assert(t, i, res, object.Integer(expd))
        case string:
            res := assertCast[object.String](t, i, obj)
            assert(t, i, res, object.String(expd))
        case Error:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, string(expd))
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

        res := assertCast[object.String](t, i, obj)
        assert(t, i, res, object.String(tst.expected))
    }
}

func TestArithmeticExpressions(t *testing.T) {
    tests := []struct {
        input    string
        expected any
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

        {"-true", UnknownOperatorError + ": -Boolean"},
        {"-true; 2", UnknownOperatorError + ": -Boolean"},
        {"true + true", UnknownOperatorError + ": Boolean + Boolean"},
        {"true + true; 2", UnknownOperatorError + ": Boolean + Boolean"},
        {`"foo" - "bar"`, UnknownOperatorError + ": String - String"},

        {"1 + true", TypeMismatchError + ": Integer + Boolean"},
        {"true + 1", TypeMismatchError + ": Boolean + Integer"},
        {"!(true + 1)", TypeMismatchError + ": Boolean + Integer"},
        {"(true + 1) * (5 + 5)", TypeMismatchError + ": Boolean + Integer"},
        {"if true + 1 { 2 }", TypeMismatchError + ": Boolean + Integer"},
        {"return true + 1", TypeMismatchError + ": Boolean + Integer"},
        {"1 + true; 2", TypeMismatchError + ": Integer + Boolean"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[object.Integer](t, i, obj)
            assert(t, i, res, object.Integer(expd))
        case Error:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, string(expd))
        }
    }
}

func TestLogicalExpressions(t *testing.T) {
    tests := []struct{
        input    string
        expected any
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

        {"!1", UnknownOperatorError + ": !Integer"},
        {"!1; 2", UnknownOperatorError + ": !Integer"},

        {"1 && 0", UnknownOperatorError + ": Integer && Integer"},
        {`"a" && "b"`, UnknownOperatorError + ": String && String"},
        {"1 || 0", UnknownOperatorError + ": Integer || Integer"},
        {`"a" || "b"`, UnknownOperatorError + ": String || String"},

        {"true && 1", TypeMismatchError + ": Boolean && Integer"},
        {"0 && false", TypeMismatchError + ": Integer && Boolean"},
        {"true || 1", TypeMismatchError + ": Boolean || Integer"},
        {"0 || false", TypeMismatchError + ": Integer || Boolean"},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case bool:
            res := assertCast[object.Boolean](t, i, obj)
            assert(t, i, res, object.Boolean(expd))
        case Error:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, string(expd))
        }
    }
}

func TestOtherErrorCases(t *testing.T) {
    tests := []struct{
        input    string
        expected string
    }{
        {"x", IdentifierNotFoundError + ": x"},
        {"!x", IdentifierNotFoundError + ": x"},
        {"if x { y }", IdentifierNotFoundError + ": x"},
        {"return x", IdentifierNotFoundError + ": x"},
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
