package parser

import (
    "fmt"
    "strconv"

    "lemur/ast"
    "lemur/lexer"
    "lemur/token"
)

const (
    _ int = iota
    Lowest
    Equals
    LessGreater
    Sum
    Product
    Prefix
    Call
)

var precedences = map[token.TokenType]int{
    token.Eq:       Equals,
    token.NotEq:    Equals,
    token.LT:       LessGreater,
    token.GT:       LessGreater,
    token.Plus:     Sum,
    token.Minus:    Sum,
    token.Slash:    Product,
    token.Asterisk: Product,
}

type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
    lex *lexer.Lexer

    errors []string
    curToken  token.Token
    nextToken token.Token

    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tt token.TokenType, f prefixParseFn) {
    p.prefixParseFns[tt] = f
}
func (p *Parser) registerInfix(tt token.TokenType, f infixParseFn) {
    p.infixParseFns[tt] = f
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        lex: l,
        errors: []string{},
    }
    p.readToken()
    p.readToken()

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.registerPrefix(token.Ident, p.parseIdentifier)
    p.registerPrefix(token.Int, p.parseIntegerLiteral)
    p.registerPrefix(token.True, p.parseBoolean)
    p.registerPrefix(token.False, p.parseBoolean)
    p.registerPrefix(token.LParen, p.parseGroupedExpression)
    p.registerPrefix(token.Bang, p.parsePrefixExpression)
    p.registerPrefix(token.Minus, p.parsePrefixExpression)
    p.registerPrefix(token.If, p.parseConditionalExpression)

    p.infixParseFns = make(map[token.TokenType]infixParseFn)
    p.registerInfix(token.Plus, p.parseInfixExpression)
    p.registerInfix(token.Minus, p.parseInfixExpression)
    p.registerInfix(token.Slash, p.parseInfixExpression)
    p.registerInfix(token.Asterisk, p.parseInfixExpression)
    p.registerInfix(token.Eq, p.parseInfixExpression)
    p.registerInfix(token.NotEq, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)

    return p
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for !p.curTokenIs(token.EOF) {
        stmt := p.parseStatement()
        p.readToken()

        if stmt == nil { continue }
        program.Statements = append(program.Statements, stmt)
    }

    return program
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) expectRead(tt token.TokenType) bool {
    if !p.nextTokenIs(tt) {
        p.expectError(tt)
        return false
    }

    p.readToken()
    return true
}

func (p *Parser) expectError(tt token.TokenType) {
    msg := fmt.Sprintf("expected %s, got %s", tt, p.nextToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(tt token.TokenType) {
    msg := fmt.Sprintf("no prefix parse function found for '%s'", tt)
    p.errors = append(p.errors, msg)
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{
        Token: p.curToken,
        Statements: []ast.Statement{},
    }

    p.readToken()
    for !p.curTokenIs(token.RBrace) {
        if p.curTokenIs(token.EOF) { return block }

        stmt := p.parseStatement()
        p.readToken()

        if stmt == nil { continue }
        block.Statements = append(block.Statements, stmt)
    }
    p.readToken()

    return block
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.Let:
        return p.parseLetStatement()
    case token.Return:
        return p.parseReturnStatement()
    case token.LBrace:
        return p.parseBlockStatement()
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}
    if !p.expectRead(token.Ident) { return nil }

    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    if !p.expectRead(token.Assign) { return nil }

    for !p.curTokenIs(token.Semicolon) { p.readToken() } // todo: parse expression here

    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}

    p.readToken()
    for !p.curTokenIs(token.Semicolon) { p.readToken() }

    return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: p.curToken}

    stmt.Value = p.parseExpression(Lowest)
    if p.nextTokenIs(token.Semicolon) { p.readToken() }

    return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]
    if prefix == nil {
        p.noPrefixParseFnError(p.curToken.Type)
        return nil
    }
    leftExp := prefix()

    for !p.nextTokenIs(token.Semicolon) && precedence < p.nextPrecedence() {
        infix := p.infixParseFns[p.nextToken.Type]
        if infix == nil { // raise error here ? (except '{', '}', ',')
            return leftExp
        }

        p.readToken()
        leftExp = infix(leftExp)
    }

    return leftExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    exp := &ast.InfixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
        Left: left,
    }

    precedence := p.curPrecedence()
    p.readToken()
    exp.Right = p.parseExpression(precedence)

    return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    exp := &ast.PrefixExpression{
        Token:    p.curToken,
        Operator: p.curToken.Literal,
    }

    p.readToken()
    exp.Right = p.parseExpression(Prefix)

    return exp
}

func (p *Parser) parseConditionalExpression() ast.Expression {
    exp:= &ast.ConditionalExpression{Token: p.curToken}

    if p.nextTokenIs(token.LParen) { p.readToken() }
    p.readToken()

    exp.Condition = p.parseExpression(Lowest)
    if p.nextTokenIs(token.RParen) { p.readToken() }

    if !p.expectRead(token.LBrace) { return nil }
    exp.Consequence = p.parseBlockStatement()

    if !p.curTokenIs(token.Else) { return exp }
    if !p.expectRead(token.LBrace) { return nil }
    exp.Alternative = p.parseBlockStatement()

    return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
    p.readToken()

    exp := p.parseExpression(Lowest)
    if !p.expectRead(token.RParen) { return nil }

    return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    l := &ast.IntegerLiteral{Token: p.curToken}

    val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
        p.errors = append(p.errors, msg)
        return nil
    }
    l.Value = val

    return l
}

func (p *Parser) parseBoolean() ast.Expression {
    return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.True)}
}

func (p *Parser) readToken() {
    p.curToken = p.nextToken
    p.nextToken = p.lex.NextToken()
}

func (p *Parser) curTokenIs(tt token.TokenType) bool { return p.curToken.Type == tt }
func (p *Parser) nextTokenIs(tt token.TokenType) bool { return p.nextToken.Type == tt }

func (p *Parser) curPrecedence() int {
    if p, ok := precedences[p.curToken.Type]; ok {
        return p
    }

    return Lowest
}

func (p *Parser) nextPrecedence() int {
    if p, ok := precedences[p.nextToken.Type]; ok {
        return p
    }

    return Lowest
}
