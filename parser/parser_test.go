package parser

import (
    "fmt"
    "testing"

    "lemur/ast"
    "lemur/lexer"
)

func TestLetStatement(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let foobar = 1729;
    `
    program := runNewParser(t, input, 3)

    tests := []struct{
        expIdentifier string
    }{
        {"x"},
        {"y"},
        {"foobar"},
    }
    for i, tst := range tests {
        stmt := program.Statements[i]
        testLetStatement(t, stmt, tst.expIdentifier)
    }
}

func TestReturnStatement(t *testing.T) {
    input := `
        return 5;
        return 10;
        return 993322;
    `
    program := runNewParser(t, input, 3)

    for _, stmt := range program.Statements {
        testReturnStatement(t, stmt)
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

        stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
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

        stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
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

    stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
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

    stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
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

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"
    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
    testIdentifier(t, stmt.Value, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"
    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
    testIntegerLiteral(t, stmt.Value, 5)
}

func TestBooleanExpression(t *testing.T) {
    input := "true;"
    program := runNewParser(t, input, 1)

    stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
    testBooleanLiteral(t, stmt.Value, true)
}

func runNewParser(t *testing.T, input string, expStatements int) *ast.Program {
    l := lexer.New(input)
    p := New(l)
    
    program := p.ParseProgram()
    checkErrors(t, p)
    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }
    if len(program.Statements) != expStatements {
        t.Fatalf("program statements should contain %d entries (got %d)",
            expStatements,
            len(program.Statements))
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

func testLetStatement(t *testing.T, s ast.Statement, expName string) {
    if s.TokenLiteral() != "let" {
        t.Errorf("statement token literal is not 'let' (got '%s')", s)
        return
    }

    ls := assertCast[*ast.LetStatement](t, s)

    if ls.Name.Value != expName {
        t.Errorf("identifier value is not '%s' (got '%s')",
            expName, ls.Name.Value)
        return
    }
    if ls.Name.TokenLiteral() != expName {
        t.Errorf("identifier token literal is not is not '%s' (got '%s')",
            expName, ls.Name.TokenLiteral())
    }
}

func testReturnStatement(t *testing.T, s ast.Statement) {
    rs := assertCast[*ast.ReturnStatement](t, s)

    if rs.TokenLiteral() != "return" {
        t.Errorf("statement token literal is not 'return' (got '%s')", rs.TokenLiteral())
    }
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, op string, right any) {
    ie := assertCast[*ast.InfixExpression](t, exp)

    if ie.Operator != op {
        t.Errorf("expression operator is not '%s' (got %s)", op, ie.Operator)
    }

    testLiteralExpression(t, ie.Left, left)
    testLiteralExpression(t, ie.Right, right)
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
    if i.TokenLiteral() != val {
        t.Errorf("identifier token literal is not %s (got %s)", val, i.TokenLiteral())
    }
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) {
    i := assertCast[*ast.IntegerLiteral](t, il)

    if i.Value != val {
        t.Errorf("integer value is not %d (got %d)", val, i.Value)
        return
    }
    if i.TokenLiteral() != fmt.Sprintf("%d", val) {
        t.Errorf("integer token literal is not %d (got %s)", val, i.TokenLiteral())
    }
}

func testBooleanLiteral(t *testing.T, be ast.Expression, val bool) {
    i := assertCast[*ast.BooleanLiteral](t, be)

    if i.Value != val {
        t.Errorf("boolean value is not %t (got %t)", val, i.Value)
        return
    }
    if i.TokenLiteral() != fmt.Sprintf("%t", val) {
        t.Errorf("boolean token literal is not %t (got %s)", val, i.TokenLiteral())
    }
}

func assertCast[T ast.Node](t *testing.T, node ast.Node) T {
    n, ok := node.(T)
    if !ok {
        t.Fatalf("node is not an %T (got %T)", *new(T), node)
    }

    return n
}

