package eval

import "github.com/threeaccents/digolang/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 || len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"isNull": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 || len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			arg := args[0]

			if arg == NULL {
				return TRUE
			}

			return FALSE
		},
	},
}
