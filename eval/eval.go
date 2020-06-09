package eval

import (
	"fmt"

	"github.com/threeaccents/digolang/ast"
	"github.com/threeaccents/digolang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(n ast.Node, env *object.Environment) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.StringLiteral:
		return &object.String{
			Value: node.Value,
		}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:

		return evalIfExpression(node, env)
	case *ast.SelectorExpression:
		left := Eval(node.Expression, env)
		if isError(left) {
			return left
		}
		return evalSelectorExpression(left, node.Selector)
	case *ast.LetStatement:
		evalLetStatement(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.FunctionLiteral:
		return &object.Function{
			Body:       node.Body,
			Parameters: node.Parameters,
			Env:        env,
		}
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}

	return nil
}

func evalSelectorExpression(left object.Object, selector *ast.Identifier) object.Object {
	if left.Type() != object.ARRAY_OBJ {
		return newError("unknown operator: %s.%s", left.Type(), selector.String())
	}

	arr := left.(*object.Array)

	obj := evalBuiltinIdentifier(selector)

	if obj.Type() != object.BUILTIN_OBJ {
		return newError("invalid selector: %s", selector.String())
	}

	builtin := obj.(*object.Builtin)

	allowedBuiltin := []string{"len", "rest", "push", "last", "first"}

	for _, item := range allowedBuiltin {
		if item == selector.String() {
			return evalBuiltin(builtin, []object.Object{arr})
		}
	}

	return newError("unknown operator: %s.%s", left.Type(), selector.String())
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	if index.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: %s%s%s", "[", index.Type(), "]")
	}

	if left.Type() != object.ARRAY_OBJ {
		return newError("unknown operator: %s%s", left.Type(), "[]")
	}

	elements := left.(*object.Array).Elements
	indexVal := index.(*object.Integer).Value

	if indexVal < 0 || indexVal > int64(len(elements)-1) {
		return NULL
	}

	return elements[indexVal]
}

func evalLetStatement(node *ast.LetStatement, env *object.Environment) {
	if node.Expression == nil {
		env.Set(node.Name.Value, NULL)
		return
	}
	val := Eval(node.Expression, env)
	if isError(val) {
		return
	}
	env.Set(node.Name.Value, val)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch funcType := fn.(type) {
	case *object.Function:
		return evalFunctionLiteral(funcType, args)
	case *object.Builtin:
		return evalBuiltin(funcType, args)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func evalBuiltin(fn *object.Builtin, args []object.Object) object.Object {
	return fn.Fn(args...)
}

func evalFunctionLiteral(fn *object.Function, args []object.Object) object.Object {
	extendedEnv := object.NewInnerEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		extendedEnv.Set(param.Value, args[paramIdx])
	}

	evaluated := Eval(fn.Body, extendedEnv)

	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return evaluated
}

func evalExpressions(arguments []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, arg := range arguments {
		v := Eval(arg, env)
		if isError(v) {
			return []object.Object{v}
		}

		result = append(result, v)
	}

	return result
}

func evalBuiltinIdentifier(node *ast.Identifier) object.Object {
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return NULL
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}

	return val
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition.Type() != object.BOOLEAN_OBJ {
		return newError("unknown operator: %s%s%s", "if(", condition.Type(), ")")
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	// since we know right and left are the same we can just pick and choose which one to use for the switch statement.
	switch left.Type() {
	case object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	}

	return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{
		Value: leftVal + rightVal,
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangPrefixExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.BOOLEAN_OBJ {
		return newError("unknown operator: %s%s", "!", right.Type())
	}

	val := right.(*object.Boolean).Value

	return nativeBoolToBooleanObject(!val)
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: %s%s", "-", right.Type())
	}

	val := right.(*object.Integer).Value

	return &object.Integer{
		Value: -val,
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func isTruthy(obj object.Object) bool {
	if obj == TRUE {
		return true
	}

	return false
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
