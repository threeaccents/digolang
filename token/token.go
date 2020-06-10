package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "hello world"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	BANG     = "!"
	GT       = ">"
	LT       = "<"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	Q

	COMMA     = ","
	PERIOD    = "."
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	RBRACKET = "]"
	LBRACKET = "["

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "if"
	RETURN   = "return"
	TRUE     = "true"
	FALSE    = "false"
	ELSE     = "else"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"else":   ELSE,
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func LookupIndentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
