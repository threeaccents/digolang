package object

import "bytes"

type Selector struct {
	Expression Object
	Selector   Object
}

func (s *Selector) Type() ObjectType { return SELECTOR_OBJ }
func (s *Selector) Inspect() string {
	var out bytes.Buffer

	out.WriteString(s.Expression.Inspect())
	out.WriteString(".")
	out.WriteString(s.Selector.Inspect())

	return out.String()
}
