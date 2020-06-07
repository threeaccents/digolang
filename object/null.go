package object

type Null struct {
}

func (i *Null) Inspect() string {
	return "null"
}

func (i *Null) Type() ObjectType {
	return NULL_OBJ
}
