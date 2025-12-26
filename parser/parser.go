// Package parser implements the parser for Nulang.
package parser

import (
	"fmt"
	"strconv"

	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/lexer"
	"github.com/nulang/nulang/token"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	COMMA       // ,
	ASSIGN      // = += -= *= /= %=
	TERNARY     // ?:
	NULLISH_C   // ??
	OR          // ||
	AND         // &&
	EQUALS      // == === != !==
	LESSGREATER // > >= < <=
	SUM         // + -
	PRODUCT     // * / %
	POWER       // **
	PREFIX      // -X !X ++X --X typeof delete
	POSTFIX     // X++ X--
	CALL        // myFunction(X)
	MEMBER      // obj.prop obj[prop]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:          ASSIGN,
	token.PLUS_ASSIGN:     ASSIGN,
	token.MINUS_ASSIGN:    ASSIGN,
	token.ASTERISK_ASSIGN: ASSIGN,
	token.SLASH_ASSIGN:    ASSIGN,
	token.PERCENT_ASSIGN:  ASSIGN,
	token.QUESTION:        TERNARY,
	token.NULLISH:         NULLISH_C,
	token.OR:              OR,
	token.AND:             AND,
	token.EQ:              EQUALS,
	token.NOT_EQ:          EQUALS,
	token.EQ3:             EQUALS,
	token.NOT_EQ3:         EQUALS,
	token.LT:              LESSGREATER,
	token.GT:              LESSGREATER,
	token.LT_EQ:           LESSGREATER,
	token.GT_EQ:           LESSGREATER,
	token.INSTANCEOF:      LESSGREATER,
	token.IN:              LESSGREATER,
	token.PLUS:            SUM,
	token.MINUS:           SUM,
	token.SLASH:           PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.PERCENT:         PRODUCT,
	token.POWER:           POWER,
	token.INCREMENT:       POSTFIX,
	token.DECREMENT:       POSTFIX,
	token.LPAREN:          CALL,
	token.LBRACKET:        MEMBER,
	token.DOT:             MEMBER,
	token.OPTIONAL:        MEMBER,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser represents the parser
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.NULL, p.parseNullLiteral)
	p.registerPrefix(token.UNDEFINED, p.parseUndefinedLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.DECREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.TYPEOF, p.parseTypeofExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseObjectLiteral)
	p.registerPrefix(token.NEW, p.parseNewExpression)
	p.registerPrefix(token.THIS, p.parseThisExpression)
	p.registerPrefix(token.ASYNC, p.parseAsyncFunction)
	p.registerPrefix(token.AWAIT, p.parseAwaitExpression)
	p.registerPrefix(token.SPREAD, p.parseSpreadExpression)
	p.registerPrefix(token.CLASS, p.parseClassExpression)
	p.registerPrefix(token.SUPER, p.parseSuperExpression)
	p.registerPrefix(token.TEMPLATE_STRING, p.parseTemplateLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.POWER, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.EQ3, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ3, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.NULLISH, p.parseInfixExpression)
	p.registerInfix(token.INSTANCEOF, p.parseInfixExpression)
	p.registerInfix(token.IN, p.parseInfixExpression)
	p.registerInfix(token.QUESTION, p.parseConditionalExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseMemberExpression)
	p.registerInfix(token.OPTIONAL, p.parseOptionalMemberExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.PLUS_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.MINUS_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.ASTERISK_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.SLASH_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.PERCENT_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerInfix(token.DECREMENT, p.parsePostfixExpression)

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors returns parse errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("line %d: expected next token to be %s, got %s instead",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("line %d: no prefix parse function for %s found",
		p.curToken.Line, t)
	p.errors = append(p.errors, msg)
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	case token.THROW:
		return p.parseThrowStatement()
	case token.TRY:
		return p.parseTryStatement()
	case token.IMPORT:
		return p.parseImportStatement()
	case token.EXPORT:
		return p.parseExportStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.CLASS:
		return p.parseClassStatement()
	case token.INTERFACE:
		return p.parseInterfaceStatement()
	case token.TYPE:
		return p.parseTypeAliasStatement()
	case token.DECLARE:
		return p.parseDeclareStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(false)
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(false)
	}

	if !p.expectPeek(token.ASSIGN) {
		p.errors = append(p.errors, "const requires initialization")
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(false)
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	// Parse initialization
	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Init = p.parseStatement()
	}

	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// Parse condition
	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	// Parse update
	if !p.curTokenIs(token.RPAREN) {
		stmt.Update = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.curToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseThrowStatement() *ast.ThrowStatement {
	stmt := &ast.ThrowStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseTryStatement() *ast.TryStatement {
	stmt := &ast.TryStatement{Token: p.curToken}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Block = p.parseBlockStatement()

	if p.peekTokenIs(token.CATCH) {
		p.nextToken()

		if p.peekTokenIs(token.LPAREN) {
			p.nextToken()
			p.nextToken()
			stmt.CatchParam = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

			if !p.expectPeek(token.RPAREN) {
				return nil
			}
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		stmt.CatchBlock = p.parseBlockStatement()
	}

	if p.peekTokenIs(token.FINALLY) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		stmt.FinallyBlock = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	p.nextToken()

	// import "module" - side-effect import
	if p.curTokenIs(token.STRING) {
		stmt.Source = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}

	// import * as name from "module"
	if p.curTokenIs(token.ASTERISK) {
		if !p.expectPeek(token.AS) {
			return nil
		}
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		stmt.NamespaceAs = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	} else if p.curTokenIs(token.LBRACE) {
		// import { a, b } from "module"
		stmt.Named = p.parseImportNames()
	} else if p.curTokenIs(token.IDENT) {
		// import name from "module"
		stmt.Default = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.FROM) {
		return nil
	}

	if !p.expectPeek(token.STRING) {
		return nil
	}

	stmt.Source = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseImportNames() []*ast.Identifier {
	names := []*ast.Identifier{}

	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return names
	}

	p.nextToken()
	names = append(names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		names = append(names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return names
}

func (p *Parser) parseExportStatement() *ast.ExportStatement {
	stmt := &ast.ExportStatement{Token: p.curToken}

	if p.peekTokenIs(token.DEFAULT) {
		p.nextToken()
		p.nextToken()
		stmt.Default = p.parseExpression(LOWEST)
	} else {
		p.nextToken()
		stmt.Named = []ast.Statement{p.parseStatement()}
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseFunctionStatement() ast.Statement {
	// Parse as expression then wrap in ExpressionStatement
	expr := p.parseFunctionLiteral()
	if expr == nil {
		return nil
	}

	funcLit, ok := expr.(*ast.FunctionLiteral)
	if ok && funcLit.Name != nil {
		// Named function declaration - treat as var declaration
		stmt := &ast.VarStatement{
			Token: funcLit.Token,
			Name:  funcLit.Name,
			Value: funcLit,
		}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}

	stmt := &ast.ExpressionStatement{Token: p.curToken, Expression: expr}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for arrow function: (x) => or x =>
	if p.peekTokenIs(token.ARROW) {
		return p.parseArrowFunction([]*ast.Identifier{ident})
	}

	return ident
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as number", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseTemplateLiteral parses template literals with ${} interpolation
func (p *Parser) parseTemplateLiteral() ast.Expression {
	tl := &ast.TemplateLiteral{Token: p.curToken}
	raw := p.curToken.Literal
	
	// Parse the raw template string into parts and expressions
	var parts []string
	var expressions []ast.Expression
	var currentPart []byte
	
	i := 0
	for i < len(raw) {
		if i < len(raw)-1 && raw[i] == '$' && raw[i+1] == '{' {
			// Save current part
			parts = append(parts, string(currentPart))
			currentPart = nil
			
			// Skip ${
			i += 2
			
			// Find matching }
			depth := 1
			exprStart := i
			for i < len(raw) && depth > 0 {
				if raw[i] == '{' {
					depth++
				} else if raw[i] == '}' {
					depth--
				}
				if depth > 0 {
					i++
				}
			}
			
			// Extract expression string and parse it
			exprStr := raw[exprStart:i]
			if len(exprStr) > 0 {
				exprLexer := lexer.New(exprStr)
				exprParser := New(exprLexer)
				expr := exprParser.parseExpression(LOWEST)
				if expr != nil {
					expressions = append(expressions, expr)
				}
			}
			i++ // Skip }
		} else {
			currentPart = append(currentPart, raw[i])
			i++
		}
	}
	
	// Add final part
	parts = append(parts, string(currentPart))
	
	tl.Parts = parts
	tl.Expressions = expressions
	
	return tl
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseUndefinedLiteral() ast.Expression {
	return &ast.UndefinedLiteral{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseTypeofExpression() ast.Expression {
	expression := &ast.TypeofExpression{Token: p.curToken}

	p.nextToken()
	expression.Value = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseAwaitExpression() ast.Expression {
	expression := &ast.AwaitExpression{Token: p.curToken}

	p.nextToken()
	expression.Value = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseSpreadExpression() ast.Expression {
	expression := &ast.SpreadExpression{Token: p.curToken}

	p.nextToken()
	expression.Value = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parseConditionalExpression(condition ast.Expression) ast.Expression {
	expression := &ast.ConditionalExpression{
		Token:     p.curToken,
		Condition: condition,
	}

	p.nextToken()
	expression.Consequence = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	expression.Alternative = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	// Check for arrow function with multiple parameters
	if p.curTokenIs(token.RPAREN) {
		// () => {}
		if p.peekTokenIs(token.ARROW) {
			return p.parseArrowFunction([]*ast.Identifier{})
		}
	}

	// Check if this could be arrow function parameters
	params := p.tryParseArrowParams()
	if params != nil {
		// Look ahead for =>, skipping potential return type
		isArrow := false
		if p.peekTokenIs(token.ARROW) {
			isArrow = true
		} else if p.peekTokenIs(token.COLON) {
			// Skip type and see if => follows
			// Since we don't want to advance the parser permanently if it's NOT an arrow,
			// this is a bit tricky. But tryParseArrowParams already advanced it to RPAREN.
			// If we see : here, it's almost certainly a type annotation for an arrow.
			isArrow = true 
		}

		if isArrow {
			return p.parseArrowFunction(params)
		}
	}

	// Regular grouped expression
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) tryParseArrowParams() []*ast.Identifier {
	// Save current state for backtracking
	savedCur := p.curToken
	savedPeek := p.peekToken
	
	params := []*ast.Identifier{}
	
	if p.curTokenIs(token.IDENT) {
		params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		
		// Handle type annotation in arrow params
		if p.peekTokenIs(token.COLON) {
			p.nextToken()
			p.skipTypeAnnotation(false)
		}
		
		for p.peekTokenIs(token.COMMA) {
			p.nextToken()
			p.nextToken()
			if !p.curTokenIs(token.IDENT) {
				// Restore state - not arrow params
				p.curToken = savedCur
				p.peekToken = savedPeek
				return nil
			}
			params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

			// Handle type annotation for subsequent arrow params
			if p.peekTokenIs(token.COLON) {
				p.nextToken()
				p.skipTypeAnnotation(false)
			}
		}
		
		if !p.peekTokenIs(token.RPAREN) {
			p.curToken = savedCur
			p.peekToken = savedPeek
			return nil
		}
		
		p.nextToken() // consume )
		return params
	}
	
	return nil
}

func (p *Parser) parseArrowFunction(params []*ast.Identifier) ast.Expression {
	fn := &ast.FunctionLiteral{
		Token:      p.curToken,
		Parameters: params,
		IsArrow:    true,
	}
	
	// Handle optional return type annotation
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(false)
	}

	if !p.expectPeek(token.ARROW) {
		return nil
	}

	p.nextToken()

	if p.curTokenIs(token.LBRACE) {
		fn.Body = p.parseBlockStatement()
	} else {
		// Expression body
		expr := p.parseExpression(LOWEST)
		fn.Body = &ast.BlockStatement{
			Statements: []ast.Statement{
				&ast.ReturnStatement{ReturnValue: expr},
			},
		}
	}

	return fn
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			nestedIf := p.parseIfExpression()
			expression.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: nestedIf},
				},
			}
		} else {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Check for async
	if p.curTokenIs(token.ASYNC) {
		lit.IsAsync = true
		p.nextToken()
	}

	// Handle named functions
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()
	
	// Handle optional return type annotation
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(true)
	}

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		lit.Body = p.parseBlockStatement()
	} else if p.peekTokenIs(token.SEMICOLON) {
		// Allowed for declarations (no body)
	} else {
		p.expectPeek(token.LBRACE) // Add error if not a semicolon or brace
		return nil
	}

	return lit
}

func (p *Parser) parseAsyncFunction() ast.Expression {
	tok := p.curToken
	p.nextToken() // move past 'async'

	if p.curTokenIs(token.FUNCTION) {
		lit := p.parseFunctionLiteral().(*ast.FunctionLiteral)
		lit.IsAsync = true
		return lit
	}

	// async arrow function: async () => {} or async x => {}
	if p.curTokenIs(token.LPAREN) || p.curTokenIs(token.IDENT) {
		var params []*ast.Identifier

		if p.curTokenIs(token.IDENT) {
			params = []*ast.Identifier{{Token: p.curToken, Value: p.curToken.Literal}}
		} else {
			p.nextToken()
			params = p.tryParseArrowParams()
		}

		fn := p.parseArrowFunction(params).(*ast.FunctionLiteral)
		fn.Token = tok
		fn.IsAsync = true
		return fn
	}

	return nil
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	// Handle optional type annotation: name: type
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(false)
	}
	
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		
		// Handle optional type annotation for subsequent parameters
		if p.peekTokenIs(token.COLON) {
			p.nextToken()
			p.skipTypeAnnotation(false)
		}
		
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{Token: p.curToken, Object: left}

	p.nextToken()
	exp.Property = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return exp
}

func (p *Parser) parseOptionalMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{Token: p.curToken, Object: left, Optional: true}

	p.nextToken()
	exp.Property = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return exp
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	exp := &ast.AssignmentExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()
	exp.Right = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.curToken}
	obj.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		
		var key ast.Expression
		var value ast.Expression
		isShorthand := false
		
		// Handle computed property names [expr]
		if p.curTokenIs(token.LBRACKET) {
			p.nextToken()
			key = p.parseExpression(LOWEST)
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		} else if p.curTokenIs(token.STRING) {
			key = p.parseStringLiteral()
		} else if p.curTokenIs(token.IDENT) {
			// Check for shorthand property: { results } => { results: results }
			identToken := p.curToken
			key = &ast.StringLiteral{Token: identToken, Value: identToken.Literal}
			
			// If next token is not COLON, it's shorthand syntax
			if !p.peekTokenIs(token.COLON) {
				isShorthand = true
				value = &ast.Identifier{Token: identToken, Value: identToken.Literal}
			}
		} else {
			key = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
		}

		if !isShorthand {
			if !p.expectPeek(token.COLON) {
				return nil
			}

			p.nextToken()
			value = p.parseExpression(LOWEST)
		}

		obj.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return obj
}

func (p *Parser) parseNewExpression() ast.Expression {
	exp := &ast.NewExpression{Token: p.curToken}

	p.nextToken()
	
	// Parse class name (with lower precedence to not include call)
	exp.Class = p.parseExpression(MEMBER)
	
	// If next token is LPAREN, parse arguments
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		exp.Arguments = p.parseExpressionList(token.RPAREN)
	}

	return exp
}

func (p *Parser) parseThisExpression() ast.Expression {
	return &ast.ThisExpression{Token: p.curToken}
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(ASSIGN))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		// Check for trailing comma - if next token is the end token, break
		if p.peekTokenIs(end) {
			p.nextToken()
			return list
		}
		p.nextToken()
		list = append(list, p.parseExpression(ASSIGN))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) isKeywordAsIdentifier() bool {
	switch p.curToken.Type {
	case token.DELETE, token.IN, token.OF, token.AS, token.FROM,
		token.GET, token.SET, token.STATIC, token.ASYNC,
		token.NEW, token.THIS, token.CLASS, token.EXTENDS,
		token.SUPER, token.TYPEOF, token.INSTANCEOF,
		token.VOID, token.YIELD, token.DEFAULT,
		token.RETURN, token.THROW, token.TRY, token.CATCH, token.FINALLY,
		token.IF, token.ELSE, token.FOR, token.WHILE, token.DO,
		token.BREAK, token.CONTINUE, token.SWITCH, token.CASE,
		token.FUNCTION, token.LET, token.CONST, token.VAR,
		token.IMPORT, token.EXPORT, token.INTERFACE, token.TYPE,
		token.PRIVATE, token.PUBLIC, token.PROTECTED, token.READONLY,
		token.IMPLEMENTS, token.NULL, token.UNDEFINED, token.TRUE, token.FALSE,
		token.DECLARE:
		return true
	}
	// Also check for "constructor" literal
	if p.curToken.Literal == "constructor" {
		return true
	}
	return false
}

// parseClassStatement parses a class declaration as a statement
func (p *Parser) parseClassStatement() ast.Statement {
	class := p.parseClassExpression()
	if class == nil {
		return nil
	}
	return &ast.ExpressionStatement{Token: class.(*ast.ClassDeclaration).Token, Expression: class}
}

// parseClassExpression parses a class declaration
func (p *Parser) parseClassExpression() ast.Expression {
	cd := &ast.ClassDeclaration{Token: p.curToken}
	
	// Class name (optional for expressions, required for declarations)
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		cd.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}
	
	// Optional extends clause
	if p.peekTokenIs(token.EXTENDS) {
		p.nextToken()
		p.nextToken()
		// Parse superclass (but stop at 'implements' or '{')
		if p.curTokenIs(token.IDENT) {
			cd.SuperClass = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			// Handle member access like stream.Readable
			for p.peekTokenIs(token.DOT) {
				p.nextToken()
				p.nextToken()
				cd.SuperClass = &ast.MemberExpression{
					Object:   cd.SuperClass,
					Property: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
				}
			}
		}
	}
	
	// Optional implements clause
	if p.peekTokenIs(token.IMPLEMENTS) {
		p.nextToken()
		cd.Implements = []*ast.Identifier{}
		
		for {
			if !p.expectPeek(token.IDENT) {
				return nil
			}
			cd.Implements = append(cd.Implements, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			
			if !p.peekTokenIs(token.COMMA) {
				break
			}
			p.nextToken()
		}
	}
	
	// Class body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	
	cd.Body = p.parseClassBody()
	
	return cd
}

// parseClassBody parses the body of a class
func (p *Parser) parseClassBody() *ast.ClassBody {
	body := &ast.ClassBody{Token: p.curToken, Members: []ast.ClassMember{}}
	
	p.nextToken()
	
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		member := p.parseClassMember()
		if member != nil {
			body.Members = append(body.Members, *member)
		}
		p.nextToken()
	}
	
	return body
}

// parseClassMember parses a class member (method or property)
func (p *Parser) parseClassMember() *ast.ClassMember {
	member := &ast.ClassMember{Token: p.curToken}
	
	// Check for modifiers: static, get, set, private, public, protected
	for {
		switch p.curToken.Type {
		case token.STATIC:
			member.IsStatic = true
			p.nextToken()
		case token.GET:
			member.IsGetter = true
			p.nextToken()
		case token.SET:
			member.IsSetter = true
			p.nextToken()
		case token.PRIVATE:
			member.IsPrivate = true
			p.nextToken()
		case token.PUBLIC, token.PROTECTED:
			// Ignore visibility modifiers
			p.nextToken()
		default:
			goto parseBody
		}
	}
	
parseBody:
	// Skip semicolons for empty statements
	if p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.RBRACE) {
		return nil
	}
	
	// Member name - allow identifiers and keywords as method/property names
	// In JavaScript, keywords can be used as method names in classes
	if p.curTokenIs(token.IDENT) || p.isKeywordAsIdentifier() {
		member.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	} else {
		return nil
	}
	
	// Check if it's a method (followed by parenthesis) or property
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(false)
	}

	if p.peekTokenIs(token.LPAREN) {
		// It's a method
		p.nextToken()
		
		// Parse parameters
		fn := &ast.FunctionLiteral{Token: p.curToken}
		fn.Parameters = p.parseFunctionParameters()

		// Handle optional return type annotation
		if p.peekTokenIs(token.COLON) {
			p.nextToken()
			p.skipTypeAnnotation(true)
		}
		
		// Parse body
		if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			fn.Body = p.parseBlockStatement()
		} else if p.peekTokenIs(token.SEMICOLON) {
			// No body
		} else {
			p.expectPeek(token.LBRACE)
			return nil
		}
		
		member.Value = fn
	} else if p.peekTokenIs(token.ASSIGN) {
		// It's a property with initializer
		p.nextToken()
		p.nextToken()
		member.Value = p.parseExpression(LOWEST)
		
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	} else if p.peekTokenIs(token.COLON) {
		// Type annotation - skip it
		p.nextToken() // skip :
		p.skipTypeAnnotation(false)
		
		if p.peekTokenIs(token.ASSIGN) {
			p.nextToken()
			p.nextToken()
			member.Value = p.parseExpression(LOWEST)
		}
		
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}
	
	return member
}

// parseSuperExpression parses 'super' keyword
func (p *Parser) parseSuperExpression() ast.Expression {
	return &ast.SuperExpression{Token: p.curToken}
}

// parseInterfaceStatement parses an interface declaration
func (p *Parser) parseInterfaceStatement() ast.Statement {
	id := &ast.InterfaceDeclaration{Token: p.curToken}
	
	// Interface name
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	id.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	// Optional extends
	if p.peekTokenIs(token.EXTENDS) {
		p.nextToken()
		id.Extends = []*ast.Identifier{}
		
		for {
			if !p.expectPeek(token.IDENT) {
				return nil
			}
			id.Extends = append(id.Extends, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			
			if !p.peekTokenIs(token.COMMA) {
				break
			}
			p.nextToken()
		}
	}
	
	// Interface body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	
	id.Body = p.parseInterfaceBody()
	
	return id
}

// parseInterfaceBody parses interface members
func (p *Parser) parseInterfaceBody() []ast.InterfaceMember {
	members := []ast.InterfaceMember{}
	
	p.nextToken()
	
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		member := p.parseInterfaceMember()
		if member != nil {
			members = append(members, *member)
		}
		p.nextToken()
	}
	
	return members
}

// parseInterfaceMember parses a single interface member
func (p *Parser) parseInterfaceMember() *ast.InterfaceMember {
	if p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.COMMA) || p.curTokenIs(token.RBRACE) {
		return nil
	}
	
	// Skip readonly modifier
	if p.curTokenIs(token.READONLY) {
		p.nextToken()
	}
	
	member := &ast.InterfaceMember{Token: p.curToken}
	
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	
	member.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	// Optional marker
	if p.peekTokenIs(token.QUESTION) {
		p.nextToken()
		member.IsOptional = true
	}
	
	// Type annotation
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.skipTypeAnnotation(false)
	}
	
	// Check for method signature
	if p.peekTokenIs(token.LPAREN) {
		member.IsMethod = true
		p.nextToken()
		// Skip parameters
		for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
			p.nextToken()
		}
		// Skip return type
		if p.peekTokenIs(token.COLON) {
			p.nextToken()
			p.skipTypeAnnotation(false)
		}
	}
	
	// Skip semicolon or comma
	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.COMMA) {
		p.nextToken()
	}
	
	return member
}

// parseDeclareStatement parses a declare statement (declare global, declare module, etc.)
func (p *Parser) parseDeclareStatement() ast.Statement {
	stmt := &ast.DeclareStatement{Token: p.curToken}
	p.nextToken() // consume 'declare'

	if p.curTokenIs(token.IDENT) {
		if p.curToken.Literal == "global" {
			p.nextToken()
			if p.curTokenIs(token.LBRACE) {
				stmt.Value = p.parseBlockStatement()
				return stmt
			}
		} else if p.curToken.Literal == "module" || p.curToken.Literal == "namespace" {
			p.nextToken()
			if p.curTokenIs(token.STRING) || p.curTokenIs(token.IDENT) {
				p.nextToken()
				if p.curTokenIs(token.LBRACE) {
					stmt.Value = p.parseBlockStatement()
					return stmt
				}
			}
		}
	}

	// Fallback to parsing as regular statement
	stmt.Value = p.parseStatement()
	
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return stmt
}

// parseTypeAliasStatement parses a type alias declaration
func (p *Parser) parseTypeAliasStatement() ast.Statement {
	ta := &ast.TypeAliasDeclaration{Token: p.curToken}
	
	// Type name
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	ta.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	// Generic parameters (skip for now)
	if p.peekTokenIs(token.LT) {
		p.skipGenericParams()
	}
	
	// Equals sign
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	
	// Skip the type definition (we don't execute types, just parse them)
	p.nextToken()
	p.skipTypeAnnotation(false)
	
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return ta
}

// skipTypeAnnotation skips over complex type annotations in a robust way
func (p *Parser) skipTypeAnnotation(stopAtBrace bool) {
	depth := 0
	for {
		switch p.peekToken.Type {
		case token.LPAREN, token.LT, token.LBRACKET:
			depth++
			p.nextToken()
		case token.LBRACE:
			if depth == 0 && stopAtBrace {
				return
			}
			depth++
			p.nextToken()
		case token.RPAREN, token.GT, token.RBRACKET, token.RBRACE:
			if depth > 0 {
				depth--
				p.nextToken()
			} else {
				// We've hit a closing delimiter at depth 0, which means the type annotation is over
				return
			}
		case token.SEMICOLON, token.COMMA, token.ASSIGN, token.EOF:
			if depth == 0 {
				// Statement or list terminator reached at depth 0
				return
			}
			p.nextToken()
		case token.ARROW:
			if depth > 0 {
				p.nextToken()
			} else {
				// Arrow function operator reached at depth 0 - stop here
				return
			}
		default:
			p.nextToken()
		}
	}
}

// skipObjectType skips an object type definition
func (p *Parser) skipObjectType() {
	depth := 1
	for depth > 0 && !p.curTokenIs(token.EOF) {
		p.nextToken()
		switch p.curToken.Type {
		case token.LBRACE:
			depth++
		case token.RBRACE:
			depth--
		}
	}
}

// skipGroupedType skips a grouped type
func (p *Parser) skipGroupedType() {
	depth := 1
	for depth > 0 && !p.curTokenIs(token.EOF) {
		p.nextToken()
		switch p.curToken.Type {
		case token.LPAREN:
			depth++
		case token.RPAREN:
			depth--
		}
	}
}

// skipBracketType skips array/tuple type notation
func (p *Parser) skipBracketType() {
	depth := 1
	for depth > 0 && !p.curTokenIs(token.EOF) {
		p.nextToken()
		switch p.curToken.Type {
		case token.LBRACKET:
			depth++
		case token.RBRACKET:
			depth--
		}
	}
}

// skipGenericParams skips generic type parameters
func (p *Parser) skipGenericParams() {
	if !p.peekTokenIs(token.LT) {
		return
	}
	p.nextToken()
	
	depth := 1
	for depth > 0 && !p.curTokenIs(token.EOF) {
		p.nextToken()
		switch p.curToken.Type {
		case token.LT:
			depth++
		case token.GT:
			depth--
		}
	}
}
