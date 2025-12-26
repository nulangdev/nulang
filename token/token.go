// Package token defines the token types for the Nulang lexer.
package token

// TokenType represents the type of a token
type TokenType string

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Token types
const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	NUMBER TokenType = "NUMBER" // 1234, 3.14
	STRING TokenType = "STRING" // "hello world"
	REGEX  TokenType = "REGEX"  // /pattern/flags

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"
	POWER    TokenType = "**"

	// Comparison
	LT     TokenType = "<"
	GT     TokenType = ">"
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT_EQ  TokenType = "<="
	GT_EQ  TokenType = ">="
	EQ3    TokenType = "==="
	NOT_EQ3 TokenType = "!=="

	// Logical
	AND         TokenType = "&&"
	OR          TokenType = "||"
	NULLISH     TokenType = "??"
	OPTIONAL    TokenType = "?."

	// Assignment operators
	PLUS_ASSIGN     TokenType = "+="
	MINUS_ASSIGN    TokenType = "-="
	ASTERISK_ASSIGN TokenType = "*="
	SLASH_ASSIGN    TokenType = "/="
	PERCENT_ASSIGN  TokenType = "%="

	// Increment/Decrement
	INCREMENT TokenType = "++"
	DECREMENT TokenType = "--"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	DOT       TokenType = "."
	QUESTION  TokenType = "?"
	ARROW     TokenType = "=>"
	SPREAD    TokenType = "..."

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// Keywords
	FUNCTION  TokenType = "FUNCTION"
	LET       TokenType = "LET"
	CONST     TokenType = "CONST"
	VAR       TokenType = "VAR"
	TRUE      TokenType = "TRUE"
	FALSE     TokenType = "FALSE"
	IF        TokenType = "IF"
	ELSE      TokenType = "ELSE"
	RETURN    TokenType = "RETURN"
	NULL      TokenType = "NULL"
	UNDEFINED TokenType = "UNDEFINED"
	FOR       TokenType = "FOR"
	WHILE     TokenType = "WHILE"
	DO        TokenType = "DO"
	BREAK     TokenType = "BREAK"
	CONTINUE  TokenType = "CONTINUE"
	SWITCH    TokenType = "SWITCH"
	CASE      TokenType = "CASE"
	DEFAULT   TokenType = "DEFAULT"
	TRY       TokenType = "TRY"
	CATCH     TokenType = "CATCH"
	FINALLY   TokenType = "FINALLY"
	THROW     TokenType = "THROW"
	NEW       TokenType = "NEW"
	THIS      TokenType = "THIS"
	CLASS     TokenType = "CLASS"
	EXTENDS   TokenType = "EXTENDS"
	SUPER     TokenType = "SUPER"
	STATIC    TokenType = "STATIC"
	IMPORT    TokenType = "IMPORT"
	EXPORT    TokenType = "EXPORT"
	FROM      TokenType = "FROM"
	AS        TokenType = "AS"
	ASYNC     TokenType = "ASYNC"
	AWAIT     TokenType = "AWAIT"
	TYPEOF    TokenType = "TYPEOF"
	INSTANCEOF TokenType = "INSTANCEOF"
	IN        TokenType = "IN"
	OF        TokenType = "OF"
	DELETE    TokenType = "DELETE"
	VOID      TokenType = "VOID"
	YIELD     TokenType = "YIELD"
	GET       TokenType = "GET"
	SET       TokenType = "SET"
	INTERFACE TokenType = "INTERFACE"
	TYPE      TokenType = "TYPE"
	IMPLEMENTS TokenType = "IMPLEMENTS"
	PRIVATE   TokenType = "PRIVATE"
	PUBLIC    TokenType = "PUBLIC"
	PROTECTED TokenType = "PROTECTED"
	READONLY  TokenType = "READONLY"
	DECLARE   TokenType = "DECLARE"
)

var keywords = map[string]TokenType{
	"function":   FUNCTION,
	"let":        LET,
	"const":      CONST,
	"var":        VAR,
	"true":       TRUE,
	"false":      FALSE,
	"if":         IF,
	"else":       ELSE,
	"return":     RETURN,
	"null":       NULL,
	"undefined":  UNDEFINED,
	"for":        FOR,
	"while":      WHILE,
	"do":         DO,
	"break":      BREAK,
	"continue":   CONTINUE,
	"switch":     SWITCH,
	"case":       CASE,
	"default":    DEFAULT,
	"try":        TRY,
	"catch":      CATCH,
	"finally":    FINALLY,
	"throw":      THROW,
	"new":        NEW,
	"this":       THIS,
	"class":      CLASS,
	"extends":    EXTENDS,
	"super":      SUPER,
	"static":     STATIC,
	"import":     IMPORT,
	"export":     EXPORT,
	"from":       FROM,
	"as":         AS,
	"async":      ASYNC,
	"await":      AWAIT,
	"typeof":     TYPEOF,
	"instanceof": INSTANCEOF,
	"in":         IN,
	"of":         OF,
	"delete":     DELETE,
	"void":       VOID,
	"yield":      YIELD,
	"get":        GET,
	"set":        SET,
	"interface":  INTERFACE,
	"type":       TYPE,
	"implements": IMPLEMENTS,
	"private":    PRIVATE,
	"public":     PUBLIC,
	"protected":  PROTECTED,
	"readonly":   READONLY,
	"declare":    DECLARE,
}

// LookupIdent checks if the identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
