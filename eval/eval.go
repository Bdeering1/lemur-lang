package eval

import (
    "fmt"

    "lemur/ast"
    "lemur/object"
)

const (
    ArgumentMistmatchError      = "wrong number of arguments for function"
    ArgumentTypesError          = "argument type(s) not supported"
    IndexOutOfBoundsError       = "index out of bounds"
    IdentifierNotFoundError     = "identifier not found"
    InfixNotImplementedError    = "no infixes implemented for type"
    InvalidConditionError       = "invalid condition"
    InvalidCastError            = "invalid type cast"
    InvalidIndexExpressionError = "invalid index expression"
    NotYetImplementedError      = "not yet implemented"
    TypeMismatchError           = "type mismatch"
    UnknownOperatorError        = "unknown operator"
    UnknownASTNodeError         = "unknown AST node"
    InternalErrorPostfix        = " (internal)"
)

var (
    True = &object.Boolean{Value: true}
    False = &object.Boolean{Value: false}
    Null = &object.Null{}
)

var builtins = map[string]object.Builtin{
    "len": func(args ...object.Object) object.Object {
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", "len")
        }

        switch input := args[0].(type) {
        case *object.Array:
            return &object.Integer{Value: int64(len(input.Elements))}
        case *object.String:
            return &object.Integer{Value: int64(len(input.Value))}
        default:
            return createError(ArgumentTypesError, "len(%s)", input.Type())
        }
    },
}

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {

    case ast.Program:
        return evalBlock(node, env)

    case *ast.BlockStatement:
        innerEnv := object.CreateEnclosedEnvironment(env)
        return evalBlock(node.Statements, innerEnv)

    case *ast.LetStatement:
        obj := Eval(node.Value, env)
        if isError(obj) { return obj }

        env.Set(node.Name.Value, obj)
        return obj

    case *ast.ReturnStatement:
        obj := Eval(node.Value, env)
        if isError(obj) { return obj }

        return &object.Return{Value: obj}

    case *ast.ExpressionStatement:
        return Eval(node.Value, env)

    case *ast.FunctionLiteral:
        return &object.Function{Parameters: node.Parameters, Body: node.Body, OuterEnv: env}

    case *ast.CallExpression:
        obj := Eval(node.Function, env)
        if isError(obj) { return obj }

        switch f := obj.(type) {
        case *object.Function:
            return evalFunction(f, node.Arguments, env)
        case object.Builtin:
            return evalBuiltin(f, node.Arguments, env)

        default:
            return createError(InvalidCastError + InternalErrorPostfix,
                "%T cannot be cast to object.Function", obj)
        }

    case *ast.ConditionalExpression:
        return evalConditionalExpression(node, env)

    case *ast.InfixExpression:
        left := Eval(node.Left, env)
        if isError(left) { return left }

        right := Eval(node.Right, env)
        if isError(right) { return right }

        return evalInfixExpression(node.Operator, left, right)

    case *ast.PrefixExpression:
        right := Eval(node.Right, env)
        if isError(right) { return right }

        return evalPrefixOperator(node.Operator, right)

    case *ast.Identifier:
        return evalIdentifier(node, env)

    case *ast.ArrayLiteral:
        arr := &object.Array{Elements: []object.Object{}}

        for _, el := range node.Elements {
            obj := Eval(el, env)
            if isError(obj) { return obj }

            arr.Elements = append(arr.Elements, obj)
        }

        return arr

    case *ast.IndexExpression:
        return evalIndexExpression(node.Left, node.Index, env)

    case *ast.StringLiteral:
        return &object.String{Value: node.Value}

    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.BooleanLiteral:
        return createBooleanObject(node.Value)

    default:
        return createError(UnknownASTNodeError + InternalErrorPostfix, "%T", node)
    }
}

func evalBlock(block []ast.Statement, env *object.Environment) object.Object {
    if len(block) == 0 { return Null } // no-op
    var obj object.Object

    for _, stmt := range block {
        obj = Eval(stmt, env)

        if obj.Type() == object.ErrorType || obj.Type() == object.ReturnType { return obj }
    }

    return obj
}

func evalBuiltin(f object.Builtin, argExprs []ast.Expression, env *object.Environment) object.Object {
    args := []object.Object{}

    for _, a := range argExprs {
        o := Eval(a, env)
        if isError(o) { return o }

        args = append(args, o)
    }

    return f(args...)
}

func evalFunction(f *object.Function, args []ast.Expression, env *object.Environment) object.Object {
    if len(args) != len(f.Parameters) {
        return createError(ArgumentMistmatchError, "%s", f)
    }

    innerEnv := object.CreateEnclosedEnvironment(f.OuterEnv)
    for i, a := range args {
        o := Eval(a, env)
        if isError(o) { return o }

        innerEnv.Set(f.Parameters[i].Value, o)
    }

    return unwrapReturn(evalBlock(f.Body.Statements, innerEnv))
}

func evalConditionalExpression(ce *ast.ConditionalExpression, env *object.Environment) object.Object {
    cond := Eval(ce.Condition, env)
    if isError(cond) { return cond }

    if cond == True { return Eval(ce.Consequence, env) }
    if cond == False {
        if ce.Alternative == nil { return Null } // default value for type or no-op
        return Eval(ce.Alternative, env)
    }

    return createError(InvalidConditionError, "%s", ce.Condition)
}

func evalIndexExpression(left, index ast.Expression, env *object.Environment) object.Object {
    leftObj := Eval(left, env)
    if isError(leftObj) { return leftObj }

    indexObj := Eval(index, env)
    if isError(indexObj) { return indexObj }

    switch {
    case leftObj.Type() == object.ArrayType && indexObj.Type() == object.IntegerType:
        arr := leftObj.(*object.Array)
        idx := indexObj.(*object.Integer).Value

        if idx < 0 || idx > int64(len(arr.Elements)) - 1 {
            return createError(IndexOutOfBoundsError, "%d", idx)
        }

        return arr.Elements[idx]

    case leftObj.Type() == object.StringType && indexObj.Type() == object.IntegerType:
        str := leftObj.(*object.String)
        idx := indexObj.(*object.Integer).Value

        if idx < 0 || idx > int64(len(str.Value)) - 1 {
            return createError(IndexOutOfBoundsError, "%d", idx)
        }

        return &object.String{Value: string(str.Value[idx])}

    default:
        return createError(
            InvalidIndexExpressionError,
            "cannot index %s with %s",
            leftObj.Type(), indexObj.Type())
    }
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
    if left.Type() != right.Type() {
        return createError(TypeMismatchError, "%s %s %s", left.Type(), operator, right.Type())
    }

    switch {
    case left.Type() == object.StringType:
        return evalStringInfixExpression(operator, left, right)
    case left.Type() == object.IntegerType:
        return evalIntegerInfixExpression(operator, left, right)
    case left.Type() == object.BooleanType:
        return evalBooleanInfixExpression(operator, left, right)
    default:
        return createError(InfixNotImplementedError, "%s", left.Type())
    }
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.String).Value
    rightVal := right.(*object.String).Value

    switch operator {
    case "+":
        return &object.String{Value: leftVal + rightVal}
    case "==":
        return createBooleanObject(leftVal == rightVal)
    case "!=":
        return createBooleanObject(leftVal != rightVal)
    default:
        return createError(UnknownOperatorError, "%s %s %s", left.Type(), operator, right.Type())
    }

}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
    case "<":
        return createBooleanObject(leftVal < rightVal)
    case ">":
        return createBooleanObject(leftVal > rightVal)
    case "==":
        return createBooleanObject(leftVal == rightVal)
    case "!=":
        return createBooleanObject(leftVal != rightVal)
    default:
        return createError(UnknownOperatorError, "%s %s %s", left.Type(), operator, right.Type())
    }
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.Boolean).Value
    rightVal := right.(*object.Boolean).Value

    switch operator {
    case "==":
        return createBooleanObject(left == right)
    case "!=":
        return createBooleanObject(left != right)
    case "&&":
        return createBooleanObject(leftVal && rightVal)
    case "||":
        return createBooleanObject(leftVal || rightVal)
    default:
        return createError(UnknownOperatorError, "%s %s %s", left.Type(), operator, right.Type())
    }
}

func evalPrefixOperator(operator string, right object.Object) object.Object {
    switch operator {
    case "!":
        return evalBangPrefix(right)
    case "-":
        return evalMinusPrefix(right)        
    default:
        return createError(UnknownOperatorError + InternalErrorPostfix, "%s%s", operator, right.Type())
    }
}

func evalBangPrefix(right object.Object) object.Object {
    switch right {
    case True:
        return False
    case False:
        return True
    default:
        return createError(UnknownOperatorError, "!%s", right.Type())
    }
}

func evalMinusPrefix(right object.Object) object.Object {
    if right.Type() != object.IntegerType {
        return createError(UnknownOperatorError, "-%s", right.Type())
    }
    
    val := right.(*object.Integer).Value
    return &object.Integer{Value: -val}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if b, ok := builtins[node.Value]; ok { return b }
    if obj, ok := env.Get(node.Value); ok { return obj }

    return createError(IdentifierNotFoundError, "%s", node.Value)
}


func unwrapReturn(obj object.Object) object.Object {
    if ret, ok := obj.(*object.Return); ok { return ret.Value }
    return obj
}

func createBooleanObject(val bool) object.Object{
    if val { return True } else { return False }
}

func createError(errKind string, msg string, args ...any) *object.Error {
    return &object.Error{Message: errKind + ": " + fmt.Sprintf(msg, args...)}
}

func isError(obj object.Object) bool { return obj.Type() == object.ErrorType }
