package parser

import (
    "lemur/ast"
    "lemur/lexer"
    "lemur/token"
)

type Parser struct {
    lex *lexer.Lexer

    curToken  token.Token
    nextToken token.Token
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{lex: l}
    p.readToken()
    p.readToken()

    return p
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.readToken()
    }

    return program
}

func (p *Parser) parseStatement() ast.Statement { // why wouldn't this be a pointer?
    return nil
}

func (p *Parser) readToken() {
    p.curToken = p.nextToken
    p.nextToken = p.lex.NextToken()
}
