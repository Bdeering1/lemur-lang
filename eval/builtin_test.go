package eval

import (
    "testing"

    "lemur/object"
)

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
        {"iter()", Error(ArgumentMistmatchError + ": iter")},
        {"iter(1)", Error(ArgumentTypesError+ ": iter(Integer)")},
        {"iter(true)", Error(ArgumentTypesError + ": iter(Boolean)")},
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)

        switch expd := tst.expected.(type) {
        case int:
            res := assertCast[object.Integer](t, i, obj)
            assert(t, i, int(res), expd)
        case []int:
            arr := assertCast[*object.Array](t, i, obj)
            for idx, el := range arr.Elements {
                res := el.(object.Integer)
                assert(t, i, int(res), expd[idx])
            }
        case Error:
            res := assertCast[*object.Error](t, i, obj)
            assert(t, i, res.Message, string(expd))
        case string:
            res := assertCast[object.String](t, i, obj)
            assert(t, i, string(res), expd)
        case nil:
            assert(t, i, obj, Null)
        }
    }
}

type IterPair struct { V any; Ok bool}

func TestBuiltinIterator(t *testing.T) {
    tests := []struct{
        input    string
        expected []IterPair
    }{
        {
            "iter([])",
            []IterPair{ {nil, false} },
        },
        {
            `iter(["one", "two"])`,
            []IterPair{ {"one", true}, {"two", true}, {"", false} },
        },
        {
            "iter([1, 2])",
            []IterPair{ {1, true}, {2, true}, {0, false} },
        },
        {
            "iter([true, false])",
            []IterPair{ {true, true}, {false, true}, {false, false} },
        },
        {
            `iter("")`,
            []IterPair{ {"", false} },
        },
        {
            `iter("abc")`,
            []IterPair{ {"a", true}, {"b", true}, {"c", true}, {"", false} },
        },
    }

    for i, tst := range tests {
        obj := runNewEval(tst.input)
        f := assertCast[object.Builtin](t, i, obj)

        for _, expd := range tst.expected {
            obj := f([]object.Object{})
            vals := assertCast[object.Tuple](t, i, obj)
            assert(t, i, len(vals), 2)

            switch expd.V.(type) {
            case string:
                val := assertCast[object.String](t, i, vals[0])
                assert(t, i, string(val), expd.V)
            case int:
                val := assertCast[object.Integer](t, i, vals[0])
                assert(t, i, int(val), expd.V)
            case bool:
                val := assertCast[object.Boolean](t, i, vals[0])
                assert(t, i, bool(val), expd.V)
            case nil:
                assert(t, i, vals[0], Null)
            }
            ok := assertCast[object.Boolean](t, i, vals[1])
            assert(t, i, bool(ok), expd.Ok)
        }
    }
}
