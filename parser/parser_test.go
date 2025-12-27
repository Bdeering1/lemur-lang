package parser

import (
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

    l := lexer.New(input)
    p := New(l)
    
    program := p.ParseProgram()
    checkErrors(t, p)
    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }
    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 entries (got %d)",
            len(program.Statements))
    }

    tests := []struct{
        expIdentifier string
    }{
        {"x"},
        {"y"},
        {"foobar"},
    }
    for i, tt := range tests {
        stmt := program.Statements[i]
        testLetStatement(t, stmt, tt.expIdentifier)
    }
}

func TestReturnStatement(t *testing.T) {
    input := `
        return 5;
        return 10;
        return 993322;
    `

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    checkErrors(t, p)
    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }
    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 entries (got %d)",
            len(program.Statements))
    }

    for _, stmt := range program.Statements {
        testReturnStatement(t, stmt)
    }
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
