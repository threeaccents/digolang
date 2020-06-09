package ast

import (
	"bytes"

	"github.com/threeaccents/digolang/token"
)

type SelectorExpression struct {
	Token      token.Token // .
	Expression Expression
	Selector   *Identifier
}

func (se *SelectorExpression) expressionNode()      {}
func (se *SelectorExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SelectorExpression) String() string {
	var out bytes.Buffer

	out.WriteString(se.Expression.String())
	out.WriteString(se.Token.Literal)
	out.WriteString(se.Selector.String())

	return out.String()
}
