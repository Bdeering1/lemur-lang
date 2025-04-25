package parser

import (
    "fmt"

    "lemur/ast"
    "lemur/lexer"
    "lemur/token"
)

type Parser struct {
    lex *lexer.Lexer

    errors []string
    curToken  token.Token
    nextToken token.Token
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        lex: l,
        errors: []string{},
    }
    p.readToken()
    p.readToken()

    return p
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for !p.curTokenIs(token.EOF) {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.readToken()
    }

    return program
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(tt token.TokenType) {
    msg := fmt.Sprintf("expected %s, got %s", tt, p.nextToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.Let:
        return p.parseLetStatement()
    case token.Return:
        return p.parseReturnStatement()
    default:
        return nil
    }
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}
    if !p.expectPeek(token.Ident) {
        return nil
    }

    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    if !p.expectPeek(token.Assign) {
        return nil
    }

    for p.curToken.Type != token.Semicolon { p.readToken() } // todo: parse expression here

    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}

    p.readToken()
    for !p.curTokenIs(token.Semicolon) { p.readToken() }

    return stmt
}

func (p *Parser) readToken() {
    p.curToken = p.nextToken
    p.nextToken = p.lex.NextToken()
}

func (p *Parser) expectPeek(tt token.TokenType) bool {
    if p.nextTokenIs(tt) {
        p.readToken()
        return true
    }
    p.peekError(tt)
    return false
}

func (p *Parser) curTokenIs(tt token.TokenType) bool { return p.curToken.Type == tt }
func (p *Parser) nextTokenIs(tt token.TokenType) bool { return p.nextToken.Type == tt }
