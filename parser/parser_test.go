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
        if !testLetStatement(t, stmt, tt.expIdentifier) { return }
    }
}

func testLetStatement(t *testing.T, s ast.Statement, expName string) bool { // is this needed?
    if s.TokenLiteral() != "let" {
        t.Errorf("s.TokenLiteral not 'let' (got %q)", s)
        return false
    }

    letStmt, ok := s.(*ast.LetStatement)
    if !ok {
        t.Errorf("s is not *ast.LetStatement (got %T)", s)
        return false
    }
    if letStmt.Name.Value != expName {
        t.Errorf("letStmt.Name.Value is not '%s' (got %q)",
            expName, letStmt.Name.Value)
        return false
    }
    if letStmt.Name.TokenLiteral() != expName {
        t.Errorf("letStmt.Name.TokenLiteral is not '%s' (got %q)",
            expName, letStmt.Name.TokenLiteral())
        return false
    }

    return true
}
