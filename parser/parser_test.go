package parser

import (
	"fmt"
	"testing"

	"lemur/ast"
	"lemur/lexer"
)

func TestLetStatement(t *testing.T) {
    tests := []struct {
        input      string
        expdIdents []string
        expdValue  any
    }{
        {"let x = 5", []string{"x"}, 5},
        {"let y = true", []string{"y"}, true},
        {"let foobar = y", []string{"foobar"}, "y"},
        {"let 1 = 2", []string{}, nil},
        {"let a, b = 1, 2", []string{}, []int{1, 2}},
        {"let a, b = true, false", []string{}, []int{1, 2}},
        {`let a, b, c = "a", "b", "c"`, []string{}, []string{"a", "b", "c"}},
    }

    for _, tst := range tests {
        switch expdVal := tst.expdValue.(type) {
        case string, int, bool:
            parser, program := runNewParser(t, tst.input)
            failOnError(t, parser)

            ls := assertCast[*ast.LetStatement](t, program[0])
            assertToken(t, ls.Token.Literal, "let")

            testIdentifier(t, ls.Names[0], tst.expdIdents[0])
            testLiteralExpression(t, ls.Value, expdVal)

        case []string:
        case []int:
        case []bool:
            parser, program := runNewParser(t, tst.input)
            failOnError(t, parser)

            ls := assertCast[*ast.LetStatement](t, program[0])
            assertToken(t, ls.Token.Literal, "let")

            exps := assertCast[ast.TupleExpression](t, ls.Value)
            assert(t, len(exps), len(expdVal))
            for i, n := range ls.Names {
                testIdentifier(t, n, (tst.expdIdents)[i])
                testLiteralExpression(t, exps[i], expdVal[i])
            }

        case nil:
            parser, _ := runNewParser(t, tst.input)
            assertError(t, parser, NonIdentifierAssignmentError)
        }
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct{
        input    string
        expValue any
    }{
        {"return 5", 5},
        {"return true", true},
        {"return y", "y"},
    }

    for _, tst := range tests {
        parser, program := runNewParser(t, tst.input)
        failOnError(t, parser)

        rs := assertCast[*ast.ReturnStatement](t, program[0])
        assertToken(t, rs.Token.Literal, "return")
        testLiteralExpression(t, rs.Value, tst.expValue)
    }
}

func TestOperatorPrecedence(t *testing.T) {
    tests := []struct{
        input    string
        expected string
    }{
        {"true", "true;"},
        {"false", "false;"},
        {"3 > 5 == false", "((3 > 5) == false);"},
        {"3 < 5 == true", "((3 < 5) == true);"},
        {"a == b && c != d", "((a == b) && (c != d));"},
        {"-a * b", "((-a) * b);"},
        {"!-a", "(!(-a));"},
        {"a + b + c", "((a + b) + c);"},
        {"a + b - c", "((a + b) - c);"},
        {"a * b * c", "((a * b) * c);"},
        {"a * b / c", "((a * b) / c);"},
        {"a + b / c", "(a + (b / c));"},
        {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f);"},
        {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4));"},
        {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4));"},
        {"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));"},
        {"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4);"},
        {"(5 + 5) * 2", "((5 + 5) * 2);"},
        {"2 / (5 + 5)", "(2 / (5 + 5));"},
        {"-(5 + 5)", "(-(5 + 5));"},
        {"!(true == true)", "(!(true == true));"},
        {"a * [1, 2, 3][b * c] * d", "((a * ([1, 2, 3][(b * c)])) * d);"},
        {"add(a[1], b * a[2], [1, 2][1] * c)", "add((a[1]), (b * (a[2])), (([1, 2][1]) * c));"},
    }

    for _, tst := range tests {
        parser, program := runNewParser(t, tst.input)
        failOnError(t, parser)

        str := program.String()
        assert(t, str, tst.expected)
    }
}

func TestTupleExpression(t *testing.T) {
    tests := []struct{
        input    string
        expected any
    }{
        {"1, 2", []int{1, 2}},
        {"1 + 1, 2 * 2", nil},
        {"true && false, a == 1", nil},
         {"fn(){}, f(1, 2)", nil},
        {`[true, true], {"a": 1, "b": 2}`, nil},
    }

    for _, tst := range tests {
        parser, program := runNewParser(t, tst.input)
        failOnError(t, parser)

        stmt := assertCast[*ast.ExpressionStatement](t, program[0])
        te := assertCast[ast.TupleExpression](t, stmt.Value)

        expd, ok := tst.expected.([]int)
        if !ok { continue }

        assert(t, len(te), len(expd))
        for i, e := range expd {
            testLiteralExpression(t, te[i], e)
        }
    }
}

func TestInfixExpression(t *testing.T) {
    infixTests := []struct{
        input    string
        leftVal  any
        operator string
        rightVal any
    }{
        {"5 + 5", 5, "+", 5},
        {"5 - 5", 5, "-", 5},
        {"5 * 5", 5, "*", 5},
        {"5 / 5", 5, "/", 5},
        {"5 > 5", 5, ">", 5},
        {"5 < 5", 5, "<", 5},
        {"5 == 5", 5, "==", 5},
        {"5 != 5", 5, "!=", 5},
        {"true == true", true, "==", true},
        {"true != false", true, "!=", false},
        {"false == false", false, "==", false},
    }

    for _, it := range infixTests {
        parser, program := runNewParser(t, it.input)
        failOnError(t, parser)

        stmt := assertCast[*ast.ExpressionStatement](t, program[0])
        testInfixExpression(t, stmt.Value, it.leftVal, it.operator, it.rightVal)
    }
}

func TestPrefixExpression(t *testing.T) {
    prefixTests := []struct{
        input    string
        operator string
        value    any
    }{
        {"!5", "!", 5},
        {"-15", "-", 15},
        {"!true", "!", true},
        {"!false", "!", false},
    }

    for _, pt := range prefixTests {
        parser, program := runNewParser(t, pt.input)
        failOnError(t, parser)

        stmt := assertCast[*ast.ExpressionStatement](t, program[0])
        exp := assertCast[*ast.PrefixExpression](t, stmt.Value)

        assertMsg(t, exp.Operator, pt.operator, "wrong expression operator")
        testLiteralExpression(t, exp.Right, pt.value)
    }
}

func TestIfExpression(t *testing.T) {
    input := "if (x < y) { x }"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    exp := assertCast[*ast.ConditionalExpression](t, stmt.Value)
    testInfixExpression(t, exp.Condition, "x", "<", "y")

    assertMsg(t, len(exp.Consequence.Statements), 1, "wrong number of statements in consequence block")
    es := assertCast[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
    testIdentifier(t, es.Value, "x")

    assertMsg(t, exp.Alternative, (*ast.BlockStatement)(nil), "wrong value for alternative block")
}

func TestIfElseExpression(t *testing.T) {
    input := "if (x < y) { x } else { y }"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    exp := assertCast[*ast.ConditionalExpression](t, stmt.Value)
    testInfixExpression(t, exp.Condition, "x", "<", "y")

    assertMsg(t, len(exp.Consequence.Statements), 1, "wrong number of statements in consequence block")
    csq := assertCast[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
    testIdentifier(t, csq.Value, "x")

    assertMsg(t, len(exp.Alternative.Statements), 1, "wrong number of statements in alternative block")
    alt := assertCast[*ast.ExpressionStatement](t, exp.Alternative.Statements[0])
    testIdentifier(t, alt.Value, "y")
}

func TestFunctionLiteral(t *testing.T) {
    input := "fn(x, y) { x + y; }"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testFunctionLiteral(t, stmt.Value, 1, []string{"x", "y"})

    f, _ := stmt.Value.(*ast.FunctionLiteral)
    es := assertCast[*ast.ExpressionStatement](t, f.Body.Statements[0])
    testInfixExpression(t, es.Value, "x", "+", "y")
}

func TestFunctionLiteralParameters(t *testing.T) {
    tests := []struct{
        input      string
        expected []string
    }{
        {input: "fn(){}", expected: []string{}},
        {input: "fn(x){}", expected: []string{"x"}},
        {input: "fn(x, y, z){}", expected: []string{"x", "y", "z"}},
    }

    for _, tst := range tests {
        parser, program := runNewParser(t, tst.input)
        failOnError(t, parser)
        stmt := assertCast[*ast.ExpressionStatement](t, program[0])

        testFunctionLiteral(t, stmt.Value, 0, tst.expected)
    }
}

func TestCallExpression(t *testing.T) {
    input := "add(1, 2 * 3, 4 + 5)"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    exp := assertCast[*ast.CallExpression](t, stmt.Value)
    testIdentifier(t, exp.Function, "add")
    assertMsg(t, len(exp.Arguments), 3, "wrong number of arguments in call expression")

    testLiteralExpression(t, exp.Arguments[0], 1)
    testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
    testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

type CallExp struct { Name string; Arguments []any }

func TestPipeOperator(t *testing.T) {
    tests := []struct{
        input     string
        expdIdent string
        expdArgs  any
    }{
        {"1 |> f", "f", []any{ 1 }},
        {"true |> f", "f", []any{ true }},
        {"a |> f", "f", []any{ "a" }},
        {"a, b |> f", "f", []any{ "a", "b" }},
        {"1 |> f1 |> f2", "f2", CallExp{ Name: "f1", Arguments: []any{ 1 } } },
    }

    for _, tst := range tests {
        parser, program := runNewParser(t, tst.input)
        failOnError(t, parser)

        stmt := assertCast[*ast.ExpressionStatement](t, program[0])
        exp := assertCast[*ast.CallExpression](t, stmt.Value)
        testIdentifier(t, exp.Function, tst.expdIdent)

        switch expd := tst.expdArgs.(type) {
        case []string:
        case []int:
        case []bool:
            assertMsg(t, len(exp.Arguments), len(expd),
                "wrong number of arguments in call expression")
            for i, expdArg := range expd {
                testLiteralExpression(t, exp.Arguments[i], expdArg)
            }
        case CallExp:
            nestedCall := assertCast[*ast.CallExpression](t, exp.Arguments[0])
            testIdentifier(t, nestedCall.Function, expd.Name)

            assertMsg(
                t, len(exp.Arguments), len(expd.Arguments),
                "wrong number of arguments in nested call expression")
            for i, a := range expd.Arguments {
                testLiteralExpression(t, nestedCall.Arguments[i], a)
            }
        }
    }
}

func TestIdentifier(t *testing.T) {
    input := "foobar;"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testIdentifier(t, stmt.Value, "foobar")
}

type KVPair struct {K any; V any}

func TestHashLiteral(t *testing.T) {
    tests := []struct{
        input    string
        expected []KVPair
    }{
        {
            "{}",
            []KVPair{},
        },
        {
            `{1: true}`,
            []KVPair{ {1, true} },
        },
        {
            `{"one": 1, "two": 2}`,
            []KVPair{ {"one", 1}, {"two", 2} },
        },
    }

    for _, tst := range tests {
        parser, program := runNewParser(t, tst.input)
        failOnError(t, parser)

        stmt := assertCast[*ast.ExpressionStatement](t, program[0])
        hl := assertCast[*ast.HashLiteral](t, stmt.Value)

        assert(t, len(hl.Pairs), len(tst.expected))
        for i, p := range hl.Pairs {
            switch expdKey := tst.expected[i].K.(type) {
            case int:
                testIntegerLiteral(t, p.Key, int64(expdKey))
            case bool:
                testBooleanLiteral(t, p.Key, expdKey)
            case string:
                testStringLiteral(t, p.Key, expdKey)
            }

            switch expdValue := tst.expected[i].V.(type) {
            case int:
                testIntegerLiteral(t, p.Value, int64(expdValue))
            case bool:
                testBooleanLiteral(t, p.Value, expdValue)
            case string:
                testStringLiteral(t, p.Value, expdValue)
            }
        }
    }
}

func TestArrayLiteral(t *testing.T) {
    input := "[1, 2 * 3, 4 + 5]"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    al := assertCast[*ast.ArrayLiteral](t, stmt.Value)
    assert(t, len(al.Elements), 3)
    testIntegerLiteral(t, al.Elements[0], 1)
    testInfixExpression(t, al.Elements[1], 2, "*", 3)
    testInfixExpression(t, al.Elements[2], 4, "+", 5)
}

func TestEmptyArrayLiteral(t *testing.T) {
    input := "[]"
    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    al := assertCast[*ast.ArrayLiteral](t, stmt.Value)
    assert(t, len(al.Elements), 0)
}

func TestIndexExpression(t *testing.T) {
    input := "arr[1 + 1]";

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    ie := assertCast[*ast.IndexExpression](t, stmt.Value)
    testIdentifier(t, ie.Left, "arr")
    testInfixExpression(t, ie.Index, 1, "+", 1)
}

func TestStringLiteral(t *testing.T) {
    input := `"foo"`;

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testStringLiteral(t, stmt.Value, "foo")
}

func TestIntegerLiteral(t *testing.T) {
    input := "5;"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testIntegerLiteral(t, stmt.Value, 5)
}

func TestBooleanExpression(t *testing.T) {
    input := "true;"

    parser, program := runNewParser(t, input)
    failOnError(t, parser)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testBooleanLiteral(t, stmt.Value, true)
}

func TestErrorCases(t *testing.T) {
    tests := []struct{
        input     string
        expdError string
    }{
        {"fn(){", EOFBeforeClosingBraceError},
        {"fn(1 + 1){}", NonIdentifierParameterError},
        {"1a", "illegal token: 1a"},
        {"let a, a = 1, 2", "duplicate identifier in variable assignment: a"},
    }

    for _, tst := range tests {
        parser, _ := runNewParser(t, tst.input)
        assertError(t, parser, tst.expdError)
    }
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, op string, right any) {
    ie := assertCast[*ast.InfixExpression](t, exp)

    assertMsg(t, ie.Operator, op, "incorrect expression operator")
    testLiteralExpression(t, ie.Left, left)
    testLiteralExpression(t, ie.Right, right)
}

func testFunctionLiteral(t *testing.T, exp ast.Expression, stmts int, params []string) {
    f := assertCast[*ast.FunctionLiteral](t, exp)

    assertMsg(t, len(f.Parameters), len(params), "wrong number of parameters in function literal")
    for i, ident := range params {
        testLiteralExpression(t, f.Parameters[i], ident)
    }

    assertMsg(t, len(f.Body.Statements), stmts, "wrong number of statements in function body")
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) { // modify to allow for strings ?
    switch v := expected.(type) {
    case int:
        testIntegerLiteral(t, exp, int64(v))
        return
    case int64:
        testIntegerLiteral(t, exp, v)
        return
    case bool:
        testBooleanLiteral(t, exp, v)
        return
    case string:
        testIdentifier(t, exp, v)
        return
    }

    t.Errorf("expression type not handled for (got %T)", exp)
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) {
    i := assertCast[*ast.Identifier](t, exp)

    assert(t, i.Value, val)
    assertToken(t, i.Token.Literal, val)
}

func testStringLiteral(t *testing.T, sl ast.Expression, val string) {
    s := assertCast[*ast.StringLiteral](t, sl)

    assert(t, s.Value, val)
    assertToken(t, s.Token.Literal, val)
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) {
    i := assertCast[*ast.IntegerLiteral](t, il)

    assert(t, i.Value, val)
    assertToken(t, i.Token.Literal, fmt.Sprintf("%d", val))
}

func testBooleanLiteral(t *testing.T, be ast.Expression, val bool) {
    i := assertCast[*ast.BooleanLiteral](t, be)

    assert(t, i.Value, val)
    assertToken(t, i.Token.Literal, fmt.Sprintf("%t", val))
}


func runNewParser(t *testing.T, input string) (*Parser, ast.Program) {
    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    if program == nil { t.Fatalf("ParseProgram() returned nil") }

    return p, program
}

func failOnError(t *testing.T, p *Parser) {
    errors := p.Errors()
    if len(errors) == 0  { return }

    for _, msg := range errors {
        t.Errorf("Error: %s", msg)
    }
    t.FailNow()
}

func assertError(t *testing.T, p *Parser, msg string) {
    errors := p.Errors()
    if len(errors) == 0 {
        t.Fatalf("no error thrown for invalid syntax")
    }

    assertMsg(t, errors[0], msg, "incorrect error message")
}

func assert(t *testing.T, val, expected any) {
    if val != expected {
        t.Fatalf("incorrect value, expected %T: %v (got %T: %v)",
            expected, expected,
            val, val)
    }
}

func assertToken(t *testing.T, val, expected any) {
    if val != expected {
        t.Errorf("incorrect token literal, expected %v (got %v)", expected, val)
    }
}

func assertMsg(t *testing.T, val, expected any, msg string) {
    if val != expected {
        t.Fatalf("%s, expected %+v (got %+v)",
            msg,
            expected, val)
    }
}

func assertCast[T ast.Node](t *testing.T, node ast.Node) T {
    n, ok := node.(T)
    if !ok {
        t.Fatalf("node is not an %T (got %T)", *new(T), node)
    }

    return n
}
