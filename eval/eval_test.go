package eval

import (
    "fmt"
    "testing"

    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
)

type Error string

func TestLetStatement(t *testing.T) {
    tests := []struct{
        input    string
        expected any
    }{
        {"let a = 5; a", 5},
        {"let a = 2 + 3; a", 5},
        {"let a = 5; let b = a; b", 5},
        {"let a = 2; let b = 3; a + b", 5},
        {"let a, b = 1, 2; a, b", []int{1, 2}},
        {"let a, b = 1 + 1, 2 * 2; a, b", []int{2, 4}},
        {"let a, b = true && false, true || false; a, b", []bool{false, true}},
        {"let a, b, c = fn(){1}(), [2][0], {3: 3}[3]; a, b, c", []int{1, 2, 3}},
        {"let a, b = fn(){ 1, 2 }()", []int{1, 2}},
        {"let a, b = fn(){ return 1, 2 }()", []int{1, 2}},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[object.Integer](t, i, obj)
            assert(t, i, res, object.Integer(expd))
        case []int:
            vals := assertCast[object.Tuple](t, i, obj)
            assert(t, i, len(vals), len(expd))

            for i, e := range expd {
                assert(t, i, vals[i], object.Integer(e))
            }
        case []bool:
            vals := assertCast[object.Tuple](t, i, obj)
            assert(t, i, len(vals), len(expd))

            for i, e := range expd {
                assert(t, i, vals[i], object.Boolean(e))
            }
        }
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct{
        input    string
        expected any
    }{
        {"return 10", 10},
        {"return 10; 9", 10},
        {"return 2 * 5; 9", 10},
        {"8; return 2 * 5; 9", 10},
        {"return 1, 2", []int{1, 2}},
        {"return 1, 2; 3", []int{1, 2}},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)
        ret := assertCast[*object.Return](t, i, obj)

        switch expd := tst.expected.(type) {
        case int:
            n := assertCast[object.Integer](t, i, ret.Value)
            assert(t, i, n, object.Integer(expd))
        case []int:
            vals := assertCast[object.Tuple](t, i, ret.Value)

            assert(t, i, len(vals), len(expd))
            for i, e := range expd {
                assert(t, i, vals[i], object.Integer(e))
            }
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

        {"if 1 + 1 { 2 }", Error(InvalidConditionError + ": (1 + 1)")},
        {`if "asdf" { 2 }`, Error(InvalidConditionError + `: "asdf"`)},
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
            assert(t, i, obj, object.Null)
        }
    }
}

type KVPair struct { K any; V any }
func TestHashLiteral(t *testing.T) {
    tests := []struct{
        input    string
        expected []KVPair
    }{
        {`{}`, []KVPair{}},
        {`{"one": "" + "1", "two": "2" + ""}`, []KVPair{ {"one", "1"}, {"two", "2"} }},
        {`{"one": 1, "two": 2}`, []KVPair{ {"one", 1}, {"two", 2} }},
        {`{"one": true, "two": false}`, []KVPair{ {"one", true} , {"two", false} }},
        {`{1: "one", 2: "two"}`, []KVPair{ {1, "one"}, {2, "two"} }},
        {`{1: 1 * 10, 2: 2 * 10}`, []KVPair{ {1, 10}, {2, 20} }},
        {`{1: true, 2: false}`, []KVPair{ {1, true}, {2, false} }},
        {`{true: "1", false: "2"}`, []KVPair{ {true, "1"}, {false, "2"} }},
        {`{true: 1, false: 2}`, []KVPair{ {true, 1}, {false, 2} }},
        {`{true: true && false, false: true || false}`, []KVPair{ {true, false}, {false, true} }},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)
        hm := assertCast[*object.HashMap](t, i, obj)
        assert(t, i, len(hm.Pairs), len(tst.expected))

        for _, p := range tst.expected {
            switch expdKey := p.K.(type) {
            case string:
                testHashLiteral(t, i, hm, object.String(expdKey), p.V)
            case int:
                testHashLiteral(t, i, hm, object.Integer(expdKey), p.V)
            case bool:
                testHashLiteral(t, i, hm, object.Boolean(expdKey), p.V)
            }
        }
    }
}

func TestInvalidHashLiteral(t *testing.T) {
    tests := []struct{
        input    string
        expected Error
    }{
        {`{[]: 0}`, Error(fmt.Sprintf(
            "%s: invalid key type (%s) for %s",
            TypeMismatchError, object.ArrayType, object.HashMapType))},
        {`{{}: 0}`, Error(fmt.Sprintf(
            "%s: invalid key type (%s) for %s",
            TypeMismatchError, object.HashMapType, object.HashMapType))},
        {`{len: 0}`, Error(fmt.Sprintf(
            "%s: invalid key type (%s) for %s",
            TypeMismatchError, object.BuiltinType, object.HashMapType))},
        {`{fn(){}: 0}`, Error(fmt.Sprintf(
            "%s: invalid key type (%s) for %s",
            TypeMismatchError, object.FunctionType, object.HashMapType))},
        {`{"one": 1, 2: 2}`, Error(fmt.Sprintf(
            "%s: %s key in %s keyed by %s",
            TypeMismatchError, object.IntegerType, object.HashMapType, object.StringType))},
        {`{1: 1, "two": 2}`, Error(fmt.Sprintf(
            "%s: %s key in %s keyed by %s",
            TypeMismatchError, object.StringType, object.HashMapType, object.IntegerType))},
        {`{"one": 1, "two": "2"}`, Error(fmt.Sprintf(
            "%s: %s value in %s with %s values",
            TypeMismatchError, object.StringType, object.HashMapType, object.IntegerType))},
        {`{"one": "1", "two": 2}`, Error(fmt.Sprintf(
            "%s: %s value in %s with %s values",
            TypeMismatchError, object.IntegerType, object.HashMapType, object.StringType))},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Error](t, i, obj)
        assert(t, i, res.Message, string(tst.expected))
    }
}

func TestArrayLiteral(t *testing.T) {
    tests := []struct{
        input    string
        expected any
    }{
        {"[1, 2]", []int{1, 2}},
        {"[true, false]", []bool{true, false}},

        {"[1, true]", Error(fmt.Sprintf(
            "%s: %s in fixed-type %s of %s",
            TypeMismatchError, object.BooleanType, object.ArrayType, object.IntegerType))},
        {"[true, 1]", Error(fmt.Sprintf(
            "%s: %s in fixed-type %s of %s",
            TypeMismatchError, object.IntegerType, object.ArrayType, object.BooleanType))},
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
        {`{"one": 1}["one"]`, 1},
        {`{"one": 1}["o" + "n" + "e"]`, 1},
        {`{1: "1"}[1]`, "1"},
        {`{1: "1"}[0 + 1]`, "1"},
        {`{true: 10}[true]`, 10},
        {`{true: 10}[true || false]`, 10},
        {`{}["asdf"]`, nil},

        {"[][0]", Error(IndexOutOfBoundsError + ": 0")},
        {"[1, 2][-1]", Error(IndexOutOfBoundsError + ": -1")},
        {`"hello"[-1]`, Error(IndexOutOfBoundsError + ": -1")},
        {`"world"[5]`, Error(IndexOutOfBoundsError + ": 5")},

        {"[1, 2][true]", Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.ArrayType, object.BooleanType))},
        {`[1, 2]["asdf"]`, Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.ArrayType, object.StringType))},
        {`""[true]`, Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.StringType, object.BooleanType))},
        {`""["asdf"]`, Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.StringType, object.StringType))},
        {`{}[[]]`, Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.HashMapType, object.ArrayType))},
        {`{}[{}]`, Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.HashMapType, object.HashMapType))},
        {`{}[len]`, Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.HashMapType, object.BuiltinType))},
        {`{}[fn(){}]`, Error(fmt.Sprintf(
            "%s: cannot index %s with %s",
            InvalidIndexExpressionError, object.HashMapType, object.FunctionType))},

        {`{"one": 1}[1]`, Error(fmt.Sprintf(
            "%s: %s key in %s keyed by %s",
            TypeMismatchError, object.IntegerType, object.HashMapType, object.StringType))},
        {`{1: "1"}["1"]`, Error(fmt.Sprintf(
            "%s: %s key in %s keyed by %s",
            TypeMismatchError, object.StringType, object.HashMapType, object.IntegerType))},
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
        case nil:
            assert(t, i, obj, object.Null)
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

        {"-true", Error(fmt.Sprintf(
            "%s: -%s",
            UnknownOperatorError, object.BooleanType))},
        {"-true; 2", Error(fmt.Sprintf(
            "%s: -%s",
            UnknownOperatorError, object.BooleanType))},
        {"true + true", Error(fmt.Sprintf(
            "%s: %s + %s",
            UnknownOperatorError, object.BooleanType, object.BooleanType))},
        {"true + true; 2", Error(fmt.Sprintf(
            "%s: %s + %s",
            UnknownOperatorError, object.BooleanType, object.BooleanType))},
        {`"foo" - "bar"`, Error(fmt.Sprintf(
            "%s: %s - %s",
            UnknownOperatorError, object.StringType, object.StringType))},

        {"1 + true", Error(fmt.Sprintf(
            "%s: %s + %s",
            TypeMismatchError, object.IntegerType, object.BooleanType))},
        {"true + 1", Error(fmt.Sprintf(
            "%s: %s + %s",
            TypeMismatchError, object.BooleanType, object.IntegerType))},
        {"!(true + 1)", Error(fmt.Sprintf(
            "%s: %s + %s",
            TypeMismatchError, object.BooleanType, object.IntegerType))},
        {"(true + 1) * (5 + 5)", Error(fmt.Sprintf(
            "%s: %s + %s",
            TypeMismatchError, object.BooleanType, object.IntegerType))},
        {"if true + 1 { 2 }", Error(fmt.Sprintf(
            "%s: %s + %s",
            TypeMismatchError, object.BooleanType, object.IntegerType))},
        {"return true + 1", Error(fmt.Sprintf(
            "%s: %s + %s",
            TypeMismatchError, object.BooleanType, object.IntegerType))},
        {"1 + true; 2", Error(fmt.Sprintf(
            "%s: %s + %s",
            TypeMismatchError, object.IntegerType, object.BooleanType))},
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

        {"!1", Error(fmt.Sprintf(
            "%s: !%s",
            UnknownOperatorError, object.IntegerType))},
        {"!1; 2", Error(fmt.Sprintf(
            "%s: !%s",
            UnknownOperatorError, object.IntegerType))},

        {"1 && 0", Error(fmt.Sprintf(
            "%s: %s && %s",
            UnknownOperatorError, object.IntegerType, object.IntegerType))},
        {`"a" && "b"`, Error(fmt.Sprintf(
            "%s: %s && %s",
            UnknownOperatorError, object.StringType, object.StringType))},
        {"1 || 0", Error(fmt.Sprintf(
            "%s: %s || %s",
            UnknownOperatorError, object.IntegerType, object.IntegerType))},
        {`"a" || "b"`, Error(fmt.Sprintf(
            "%s: %s || %s",
            UnknownOperatorError, object.StringType, object.StringType))},

        {"true && 1", Error(fmt.Sprintf(
            "%s: %s && %s",
            TypeMismatchError, object.BooleanType, object.IntegerType))},
        {"0 && false", Error(fmt.Sprintf(
            "%s: %s && %s",
            TypeMismatchError, object.IntegerType, object.BooleanType))},
        {"true || 1", Error(fmt.Sprintf(
            "%s: %s || %s",
            TypeMismatchError, object.BooleanType, object.IntegerType))},
        {"0 || false", Error(fmt.Sprintf(
            "%s: %s || %s",
            TypeMismatchError, object.IntegerType, object.BooleanType))},
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
        expected Error
    }{
        {"x", Error(IdentifierNotFoundError + ": x")},
        {"x", Error(IdentifierNotFoundError + ": x")},
        {"!x", Error(IdentifierNotFoundError + ": x")},
        {"if x { y }", Error(IdentifierNotFoundError + ": x")},
        {"return x", Error(IdentifierNotFoundError + ": x")},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        res := assertCast[*object.Error](t, i, obj)
        assert(t, i, res.Message, string(tst.expected))
    }
}


func testHashLiteral[K object.Object](t *testing.T, i int, hm *object.HashMap, ek K, ev any) {
    switch ev := ev.(type) {
    case string:
        res := assertCast[object.String](t, i, hm.Pairs[ek])
        assert(t, i, res, object.String(ev))
    case int:
        res := assertCast[object.Integer](t, i, hm.Pairs[ek])
        assert(t, i, res, object.Integer(ev))
    case bool:
        res := assertCast[object.Boolean](t, i, hm.Pairs[ek])
        assert(t, i, res, object.Boolean(ev))
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
