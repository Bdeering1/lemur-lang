package parser

import (
    "fmt"
    "testing"

    "lemur/ast"
    "lemur/lexer"
)

func TestLetStatement(t *testing.T) {
    tests := []struct {
        input         string
        expIdentifier string
        expValue      any
    }{
        {"let x = 5;", "x", 5},
        {"let y = true;", "y", true},
        {"let foobar = y;", "foobar", "y"},
    }

    for _, tst := range tests {
        program := runNewParser(t, tst.input, 1)
        ls := assertCast[*ast.LetStatement](t, program[0])

        if ls.Token.Literal != "let" {
            t.Errorf("statement token literal is not 'let' (got '%s')", ls)
            return
        }

        testIdentifier(t, ls.Name, tst.expIdentifier)
        testLiteralExpression(t, ls.Value, tst.expValue)
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct{
        input    string
        expValue any
    }{
        {"return 5;", 5},
        {"return true;", true},
        {"return y;", "y"},
    }

    for _, tst := range tests {
        program := runNewParser(t, tst.input, 1)
        rs := assertCast[*ast.ReturnStatement](t, program[0])

        if rs.Token.Literal != "return" {
            t.Errorf("statement token literal is not 'return' (got '%s')", rs.Token.Literal)
        }
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
    }

    for _, tst := range tests {
        program := runNewParser(t, tst.input, 1)

        str := program.String()
        if str != tst.expected {
            t.Errorf("expected %q (got %q)", tst.expected, str)
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
        program := runNewParser(t, it.input, 1)

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
        program := runNewParser(t, pt.input, 1)

        stmt := assertCast[*ast.ExpressionStatement](t, program[0])
        exp := assertCast[*ast.PrefixExpression](t, stmt.Value)

        if exp.Operator != pt.operator {
            t.Errorf("expression operator is not '%s' (got %s)",
                pt.operator,
                exp.Operator)
        }
        testLiteralExpression(t, exp.Right, pt.value)
    }
}

func TestIfExpression(t *testing.T) {
    input := "if (x < y) { x }"

    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    exp := assertCast[*ast.ConditionalExpression](t, stmt.Value)
    testInfixExpression(t, exp.Condition, "x", "<", "y")

    if (len(exp.Consequence.Statements) != 1) {
        t.Errorf("consequence statements should contain 1 entry (got %d)",
            len(exp.Consequence.Statements))
    }
    es := assertCast[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
    testIdentifier(t, es.Value, "x")

    if (exp.Alternative != nil) {
        t.Errorf("conditional expression alternative should be nil (got %+v)", exp.Alternative)
    }
}

func TestIfElseExpression(t *testing.T) {
    input := "if (x < y) { x } else { y }"

    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    exp := assertCast[*ast.ConditionalExpression](t, stmt.Value)
    testInfixExpression(t, exp.Condition, "x", "<", "y")

    if (len(exp.Consequence.Statements) != 1) {
        t.Errorf("consequence statements should contain 1 entry (got %d)",
            len(exp.Consequence.Statements))
    }
    csq := assertCast[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
    testIdentifier(t, csq.Value, "x")

    if (len(exp.Alternative.Statements) != 1) {
        t.Errorf("alternative statements should contain one entry (got %d)",
            len(exp.Alternative.Statements))
    }
    alt := assertCast[*ast.ExpressionStatement](t, exp.Alternative.Statements[0])
    testIdentifier(t, alt.Value, "y")
}

func TestFunctionLiteral(t *testing.T) {
    input := "fn(x, y) { x + y; }"

    program := runNewParser(t, input, 1)
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
        program := runNewParser(t, tst.input, 1)
        stmt := assertCast[*ast.ExpressionStatement](t, program[0])

        testFunctionLiteral(t, stmt.Value, 0, tst.expected)
    }
}

func TestCallExpression(t *testing.T) {
    input := "add(1, 2 * 3, 4 + 5)"

    program := runNewParser(t, input, 1)
    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    exp := assertCast[*ast.CallExpression](t, stmt.Value)

    testIdentifier(t, exp.Function, "add")

    if len(exp.Arguments) != 3 {
        t.Errorf("call expression should contain %d arguments (got %d)", 3, exp.Arguments)
    }
    testLiteralExpression(t, exp.Arguments[0], 1)
    testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
    testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"
    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testIdentifier(t, stmt.Value, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"
    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testIntegerLiteral(t, stmt.Value, 5)
}

func TestBooleanExpression(t *testing.T) {
    input := "true;"
    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program[0])
    testBooleanLiteral(t, stmt.Value, true)
}

func runNewParser(t *testing.T, input string, expStatements int) ast.Program {
    l := lexer.New(input)
    p := New(l)
    
    program := p.ParseProgram()
    checkErrors(t, p)
    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }
    if len(program) != expStatements {
        t.Fatalf("program statements should contain %d entries (got %d)",
            expStatements,
            len(program))
    }

    return program
}

func checkErrors(t *testing.T, p *Parser) {
    errors := p.Errors()
    if len(errors) == 0  { return }

    t.Errorf("parser has %d errors:", len(errors))
    for _, msg := range errors {
        t.Errorf("\t%q", msg)
    }
    t.FailNow()
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, op string, right any) {
    ie := assertCast[*ast.InfixExpression](t, exp)

    if ie.Operator != op {
        t.Errorf("expression operator is not '%s' (got %s)", op, ie.Operator)
    }

    testLiteralExpression(t, ie.Left, left)
    testLiteralExpression(t, ie.Right, right)
}

func testFunctionLiteral(t *testing.T, exp ast.Expression, stmts int, params []string) {
    f := assertCast[*ast.FunctionLiteral](t, exp)

    if len(f.Parameters) != len(params) {
        t.Fatalf("wrong number of parameters in function literal, should be %d (got %d)",
            len(params),
            len(f.Parameters))
    }

    for i, ident := range params {
        testLiteralExpression(t, f.Parameters[i], ident)
    }

    if len(f.Body.Statements) != stmts {
        t.Fatalf("function body should contain %d statements (got %d)",
            len(f.Body.Statements),
            stmts)
    }
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) {
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

    if i.Value != val {
        t.Errorf("identifier value is not %s (got %s)", val, i.Value)
    }
    if i.Token.Literal != val {
        t.Errorf("identifier token literal is not %s (got %s)", val, i.Token.Literal)
    }
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) {
    i := assertCast[*ast.IntegerLiteral](t, il)

    if i.Value != val {
        t.Errorf("integer value is not %d (got %d)", val, i.Value)
        return
    }
    if i.Token.Literal != fmt.Sprintf("%d", val) {
        t.Errorf("integer token literal is not %d (got %s)", val, i.Token.Literal)
    }
}

func testBooleanLiteral(t *testing.T, be ast.Expression, val bool) {
    i := assertCast[*ast.BooleanLiteral](t, be)

    if i.Value != val {
        t.Errorf("boolean value is not %t (got %t)", val, i.Value)
        return
    }
    if i.Token.Literal != fmt.Sprintf("%t", val) {
        t.Errorf("boolean token literal is not %t (got %s)", val, i.Token.Literal)
    }
}

func assertCast[T ast.Node](t *testing.T, node ast.Node) T {
    n, ok := node.(T)
    if !ok {
        t.Fatalf("node is not an %T (got %T)", *new(T), node)
    }

    return n
}
