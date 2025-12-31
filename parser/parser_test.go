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

func TestPrefixExpression(t *testing.T) {
    prefixTests := []struct{
        input    string
        operator string
        intValue int64
    }{
        {"!5;", "!", 5},
        {"-15;", "-", 15},
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
        testIntegerLiteral(t, exp.Right, pt.intValue)
    }
}

func TestInfixExpression(t *testing.T) {
    infixTests := []struct{
        input    string
        leftVal  int64
        operator string
        rightVal int64
    }{
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},
    }

    for _, it := range infixTests {
        program := runNewParser(t, it.input, 1)

        stmt := assertCast[*ast.ExpressionStatement](t, program.Statements[0])
        testInfixExpression(t, stmt.Value, it.leftVal, it.operator, it.rightVal)
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct{
        input    string
        expected string
    }{
        {"-a * b", "((-a) * b);"},
        {"!-a", "(!(-a));"},
        {"a + b + c", "((a + b) + c);"},
        {"a + b - c", "((a + b) - c);"},
        {"a * b * c", "((a * b) * c);"},
        {"a * b / c", "((a * b) / c);"},
        {"a + b / c", "(a + (b / c));"},
        {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f);"},
        // {"3 + 4; -5 * 5", "(3 + 4);((-5) * 5);"},
        {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4));"},
        {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4));"},
        {"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));"},
    }

    for _, tst := range tests {
        program := runNewParser(t, tst.input, 1)

        str := program.String()
        if str != tst.expected {
            t.Errorf("expected %q (got %q)", tst.expected, str)
        }
    }
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
        t.Fatalf("program statements does not contain %d entries (got %d)",
            expStatements,
            len(program.Statements))
    }

    return program
}

func assertCast[T ast.Node](t *testing.T, node ast.Node) T {
    n, ok := node.(T)
    if !ok {
        t.Fatalf("node is not an %T (got %T)", *new(T), node)
    }

    return n
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
        t.Errorf("s.TokenLiteral not 'let' (got '%s')", s)
        return
    }

    letStmt, ok := s.(*ast.LetStatement)
    if !ok {
        t.Errorf("s is not *ast.LetStatement (got %T)", s)
        return
    }
    if letStmt.Name.Value != expName {
        t.Errorf("letStmt.Name.Value is not '%s' (got '%s')",
            expName, letStmt.Name.Value)
        return
    }
    if letStmt.Name.TokenLiteral() != expName {
        t.Errorf("letStmt.Name.TokenLiteral is not '%s' (got '%s')",
            expName, letStmt.Name.TokenLiteral())
    }
}

func testReturnStatement(t *testing.T, s ast.Statement) {
    returnStmt, ok := s.(*ast.ReturnStatement)
    if !ok {
        t.Errorf("s is not *ast.ReturnStatement (got %T)", s)
        return
    }
    if returnStmt.TokenLiteral() != "return" {
        t.Errorf("s.TokenLiteral not 'return' (got '%s')", returnStmt.TokenLiteral())
    }
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, op string, right any) {
    ie, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Errorf("expression '%s' is not an ast.InfixExpression (got %T)", exp, exp)
    }

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
    case string:
        testIdentifier(t, exp, v)
        return
    }

    t.Errorf("expression type not handled (got %T)", exp)
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) {
    i, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il is not an *ast.IntegerLiteral (got %T)", il)
        return
    }

    if i.Value != val {
        t.Errorf("i.Value is not %d (got %d)", val, i.Value)
        return
    }

    if i.TokenLiteral() != fmt.Sprintf("%d", val) {
        t.Errorf("i.TokenLiteral is not %d (got %s)", val, i.TokenLiteral())
    }
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) {
    ident, ok := exp.(*ast.Identifier)
    if !ok {
        t.Fatalf("expression is not an ast.Identifier (got %T)", exp)
    }
    if ident.Value != val {
        t.Errorf("ident.Value is not %s (got %s)", val, ident.Value)
    }
    if ident.TokenLiteral() != val {
        t.Errorf("ident.TokenLiteral not %s (got %s)", val, ident.TokenLiteral())
    }
}
