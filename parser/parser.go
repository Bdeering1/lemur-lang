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
    token.LParen:   Call,
}

type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
    lex *lexer.Lexer

    errors  []string
    invalid   bool
    curToken  token.Token

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

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.registerPrefix(token.Ident, p.parseIdentifier)
    p.registerPrefix(token.String, p.parseStringLiteral)
    p.registerPrefix(token.Int, p.parseIntegerLiteral)
    p.registerPrefix(token.True, p.parseBoolean)
    p.registerPrefix(token.False, p.parseBoolean)
    p.registerPrefix(token.LParen, p.parseGroupedExpression)
    p.registerPrefix(token.Bang, p.parsePrefixOperator)
    p.registerPrefix(token.Minus, p.parsePrefixOperator)
    p.registerPrefix(token.If, p.parseConditionalExpression)
    p.registerPrefix(token.Function, p.parseFunctionLiteral)

    p.infixParseFns = make(map[token.TokenType]infixParseFn)
    p.registerInfix(token.Plus, p.parseInfixExpression)
    p.registerInfix(token.Minus, p.parseInfixExpression)
    p.registerInfix(token.Slash, p.parseInfixExpression)
    p.registerInfix(token.Asterisk, p.parseInfixExpression)
    p.registerInfix(token.Eq, p.parseInfixExpression)
    p.registerInfix(token.NotEq, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
    p.registerInfix(token.LParen, p.parseCallExpression)

    return p
}

func (p *Parser) ParseProgram() ast.Program {
    program := ast.Program{}

    for !p.curTokenIs(token.EOF) {
        stmt := p.parseStatement()

        if p.invalid { break }
        program = append(program, stmt)
    }

    return program
}

func (p *Parser) Errors() []string {
    return p.errors
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{
        Token: p.curToken,
        Statements: []ast.Statement{},
    }
    p.readToken()

    for !p.curTokenIs(token.RBrace) {
        if p.curTokenIs(token.EOF) {
            p.raiseError("Reach EOF before closing brace in block statement (missing '}')")
            return block
        }

        stmt := p.parseStatement()

        if stmt == nil { continue }
        block.Statements = append(block.Statements, stmt)
    }
    p.readToken()

    return block
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}
    p.readToken()

    if !p.curTokenIs(token.Ident) { return nil }
    stmt.Name, _ = p.parseIdentifier().(*ast.Identifier)

    if !p.expectRead(token.Assign) { return nil }
    stmt.Value = p.parseExpression(Lowest)

    if p.curTokenIs(token.Semicolon) { p.readToken() }

    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}
    p.readToken()

    stmt.Value = p.parseExpression(Lowest)
    if p.curTokenIs(token.Semicolon) { p.readToken() }

    return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: p.curToken}

    stmt.Value = p.parseExpression(Lowest)
    if stmt.Value == nil { return nil }

    if p.curTokenIs(token.Semicolon) { p.readToken() }

    return stmt
}


func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]
    if prefix == nil {
        p.noPrefixParseFnError()
        return nil
    }

    exp := prefix()
    if exp == nil { return nil }

    for !p.curTokenIs(token.Semicolon) && precedence < p.curPrecedence() {
        infix := p.infixParseFns[p.curToken.Type]
        if infix == nil {
            p.noInfixParseFnError()
            return exp
        }

        exp = infix(exp)
    }

    return exp
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

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    exp := &ast.CallExpression{Token: p.curToken, Function: function}
    p.readToken()

    exp.Arguments = p.parseCallArguments()

    return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if p.skipToken(token.RParen) { return args }

    args = append(args, p.parseExpression(Lowest))
    for p.curTokenIs(token.Comma) {
        p.readToken()
        args = append(args, p.parseExpression(Lowest))
    }

    if !p.expectRead(token.RParen) { return nil }
    return args
}

func (p *Parser) parsePrefixOperator() ast.Expression {
    exp := &ast.PrefixExpression{
        Token:    p.curToken,
        Operator: p.curToken.Literal,
    }
    p.readToken()

    exp.Right = p.parseExpression(Prefix)

    return exp
}

func (p *Parser) parseConditionalExpression() ast.Expression {
    exp := &ast.ConditionalExpression{Token: p.curToken}
    p.readToken()

    if p.curTokenIs(token.LParen) { p.readToken() }

    exp.Condition = p.parseExpression(Lowest)
    if p.curTokenIs(token.RParen) { p.readToken() }

    if !p.curTokenIs(token.LBrace) { return nil }
    exp.Consequence = p.parseBlockStatement()

    if !p.curTokenIs(token.Else) { return exp }
    p.readToken()

    if !p.curTokenIs(token.LBrace) { return nil }
    exp.Alternative = p.parseBlockStatement()

    return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
    l := &ast.FunctionLiteral{Token: p.curToken}
    p.readToken()

    if !p.expectRead(token.LParen) { return nil }
    l.Parameters = p.parseFunctionParameters()

    if l.Parameters == nil || !p.curTokenIs(token.LBrace) { return nil }
    l.Body = p.parseBlockStatement()

    return l
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
    idents := []*ast.Identifier{}

    if p.skipToken(token.RParen) { return idents }
    i, _ := p.parseIdentifier().(*ast.Identifier)
    idents = append(idents, i)

    for p.curTokenIs(token.Comma) {
        p.readToken()
        i, _ := p.parseIdentifier().(*ast.Identifier)
        idents = append(idents, i)
    }

    if !p.expectRead(token.RParen) { return nil }

    return idents
}

func (p *Parser) parseGroupedExpression() ast.Expression {
    p.readToken()

    exp := p.parseExpression(Lowest)
    if !p.expectRead(token.RParen) { return nil }

    return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
    i := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    p.readToken()

    return i
}

func (p *Parser) parseStringLiteral() ast.Expression {
    s := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
    p.readToken()

    return s
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

    p.readToken()
    return l
}

func (p *Parser) parseBoolean() ast.Expression {
    b := &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.True)}
    p.readToken()

    return b
}

func (p *Parser) readToken() { p.curToken = p.lex.NextToken() }
func (p *Parser) curTokenIs(tt token.TokenType) bool { return p.curToken.Type == tt }

func (p *Parser) curPrecedence() int {
    if p, ok := precedences[p.curToken.Type]; ok {
        return p
    }

    return Lowest
}

func (p *Parser) expectRead(tt token.TokenType) bool {
    if !p.curTokenIs(tt) {
        p.expectError(tt)
        return false
    }

    p.readToken()
    return true
}

func (p *Parser) skipToken(tt token.TokenType) bool {
    if !p.curTokenIs(tt) {
        return false
    }

    p.readToken()
    return true
}


func (p *Parser) expectError(tt token.TokenType) {
    p.raiseError(fmt.Sprintf("expected %s, got %s", tt, p.curToken.Type))
}

func (p *Parser) noPrefixParseFnError() {
    p.raiseError(fmt.Sprintf("no prefix parse function found for '%s'", p.curToken.Type))
}

func (p *Parser) noInfixParseFnError() {
    p.raiseError(fmt.Sprintf("no infix parse function found for '%s'", p.curToken.Type))
}

func (p *Parser) raiseError(msg string) {
    p.errors = append(p.errors, msg)
    p.invalid = true
}
