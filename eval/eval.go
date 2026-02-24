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
    NotEnoughValuesError           = "not enough values in assignment"
    NotYetImplementedError      = "not yet implemented"
    TypeMismatchError           = "type mismatch"
    UnknownOperatorError        = "unknown operator"
    UnknownASTNodeError         = "unknown AST node"
    InternalErrorPostfix        = " (internal)"
)

var b BuiltinProvider = B{}
var builtins map[string]object.Builtin
func init() { builtins = b.Builtins() }


func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {

    case ast.Program:
        return evalBlock(node, env)

    case *ast.BlockStatement:
        innerEnv := object.CreateEnclosedEnvironment(env)
        return evalBlock(node.Statements, innerEnv)

    case *ast.LetStatement:
        obj := Eval(node.Values, env)
        if isError(obj) { return obj }

        values := obj.(object.Tuple)
        if len(values) != len(node.Names) {
            return createError(
                NotEnoughValuesError,
                "expected %d values (got %d)",
                len(node.Names), len(values))
        }

        for i, n := range node.Names  {
            env.Set(n.Value, values[i])
        }
        return values

    case *ast.ReturnStatement:
        obj := Eval(node.Value, env)
        if isError(obj) { return obj }

        return &object.Return{Value: obj}

    case *ast.ExpressionStatement:
        return Eval(node.Value, env)

    case ast.TupleExpression:
        return evalTupleExpression(node, env)

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
            return createError(
                InvalidCastError + InternalErrorPostfix,
                "%T cannot be cast to object.Function",
                obj)
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

    case *ast.HashLiteral:
        return evalHashLiteral(node, env)

    case *ast.ArrayLiteral:
        return evalArray(node, env)

    case *ast.IndexExpression:
        return evalIndexExpression(node.Left, node.Index, env)

    case *ast.StringLiteral:
        return object.String(node.Value)

    case *ast.IntegerLiteral:
        return object.Integer(node.Value)

    case *ast.BooleanLiteral:
        return object.Boolean(node.Value)

    default:
        return createError(UnknownASTNodeError + InternalErrorPostfix, "%T", node)
    }
}

func evalBlock(block []ast.Statement, env *object.Environment) object.Object {
    if len(block) == 0 { return object.Null } // no-op
    var obj object.Object

    for _, stmt := range block {
        obj = Eval(stmt, env)

        if obj.Type() == object.ErrorType || obj.Type() == object.ReturnType { return obj }
    }

    return obj
}

func evalTupleExpression(tuple ast.TupleExpression, env *object.Environment) object.Object {
    t := object.Tuple{}

    for _, e := range tuple {
        o := Eval(e, env)
        if isError(o) { return o }

        t = append(t, o)
    }

    return t
}

func evalBuiltin(f object.Builtin, argExprs []ast.Expression, env *object.Environment) object.Object {
    args := []object.Object{}

    for _, a := range argExprs {
        o := Eval(a, env)
        if isError(o) { return o }

        args = append(args, o)
    }

    return f(args, env)
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
    obj := Eval(ce.Condition, env)
    if isError(obj) { return obj }

    cond, ok := obj.(object.Boolean)
    if !ok {
        return createError(InvalidConditionError, "%s", ce.Condition)
    }

    if cond { return Eval(ce.Consequence, env) }

    if ce.Alternative == nil { return object.Null } // default value for type or no-op
    return Eval(ce.Alternative, env)
}

func evalIndexExpression(left, index ast.Expression, env *object.Environment) object.Object {
    leftObj := Eval(left, env)
    if isError(leftObj) { return leftObj }

    indexObj := Eval(index, env)
    if isError(indexObj) { return indexObj }

    switch {
    case leftObj.Type() == object.HashMapType && isValueObject(indexObj):
        hm := leftObj.(*object.HashMap)

        if hm.KeyType == object.NullType { return object.Null }
        if indexObj.Type() != hm.KeyType {
            return createError(
                TypeMismatchError,
                "%s key in %s keyed by %s",
                indexObj.Type(), object.HashMapType, hm.KeyType)
        }

        res, ok := hm.Pairs[indexObj]
        if !ok { return object.Null }

        return res

    case leftObj.Type() == object.ArrayType && indexObj.Type() == object.IntegerType:
        arr := leftObj.(*object.Array)
        idx := int(indexObj.(object.Integer))

        if idx < 0 || idx > len(arr.Elements) - 1 {
            return createError(IndexOutOfBoundsError, "%d", idx)
        }

        return arr.Elements[idx]

    case leftObj.Type() == object.StringType && indexObj.Type() == object.IntegerType:
        str := leftObj.(object.String)
        idx := int(indexObj.(object.Integer))

        if idx < 0 || idx > len(str) - 1 {
            return createError(IndexOutOfBoundsError, "%d", idx)
        }

        return object.String(str[idx])

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
    leftVal := left.(object.String)
    rightVal := right.(object.String)

    switch operator {
    case "+":
        return object.String(leftVal + rightVal)
    case "==":
        return object.Boolean(leftVal == rightVal)
    case "!=":
        return object.Boolean(leftVal != rightVal)
    default:
        return createError(UnknownOperatorError, "%s %s %s", left.Type(), operator, right.Type())
    }

}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(object.Integer)
    rightVal := right.(object.Integer)

    switch operator {
    case "+":
        return object.Integer(leftVal + rightVal)
    case "-":
        return object.Integer(leftVal - rightVal)
    case "*":
        return object.Integer(leftVal * rightVal)
    case "/":
        return object.Integer(leftVal / rightVal)
    case "<":
        return object.Boolean(leftVal < rightVal)
    case ">":
        return object.Boolean(leftVal > rightVal)
    case "==":
        return object.Boolean(leftVal == rightVal)
    case "!=":
        return object.Boolean(leftVal != rightVal)
    default:
        return createError(UnknownOperatorError, "%s %s %s", left.Type(), operator, right.Type())
    }
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(object.Boolean)
    rightVal := right.(object.Boolean)

    switch operator {
    case "==":
        return object.Boolean(left == right)
    case "!=":
        return object.Boolean(left != right)
    case "&&":
        return leftVal && rightVal
    case "||":
        return leftVal || rightVal
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
    b, ok := right.(object.Boolean)
    if !ok {
        return createError(UnknownOperatorError, "!%s", right.Type())
    }

    return !b
}

func evalMinusPrefix(right object.Object) object.Object {
    if right.Type() != object.IntegerType {
        return createError(UnknownOperatorError, "-%s", right.Type())
    }
    
    val := right.(object.Integer)
    return object.Integer(-val)
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if b, ok := builtins[node.Value]; ok { return b }
    if obj, ok := env.Get(node.Value); ok { return obj }

    return createError(IdentifierNotFoundError, "%s", node.Value)
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
    hm := &object.HashMap{
        Pairs: map[object.Object]object.Object{},
        KeyType: object.NullType,
        ValueType: object.NullType,
    }
    if len(node.Pairs) == 0 { return hm }

    key := Eval(node.Pairs[0].Key, env)
    if isError(key) { return key }

    if !isValueObject(key) {
        return createError(
            TypeMismatchError,
            "invalid key type (%s) for %s",
            key.Type(), object.HashMapType)
    }
    hm.KeyType = key.Type()

    val := Eval(node.Pairs[0].Value, env)
    if isError(val) { return val }
    hm.ValueType = val.Type()

    for _, p := range node.Pairs {
        key := Eval(p.Key, env)
        if isError(key) { return key }
        if key.Type() != hm.KeyType {
            return createError(
                TypeMismatchError,
                "%s key in %s keyed by %s",
                key.Type(), object.HashMapType, hm.KeyType)
        }

        val := Eval(p.Value, env)
        if isError(val) { return val }
        if val.Type() != hm.ValueType {
            return createError(
                TypeMismatchError,
                "%s value in %s with %s values",
                val.Type(), object.HashMapType, hm.ValueType)
        }

        hm.Pairs[key] = val
    }

    return hm
}

func evalArray(node *ast.ArrayLiteral, env *object.Environment) object.Object {
    arr := &object.Array{
        Elements: []object.Object{},
        ElementType: object.NullType,
    }
    if len(node.Elements) == 0 { return arr }

    o := Eval(node.Elements[0], env)
    if isError(o) { return o }

    arr.ElementType = o.Type()
    arr.Elements = append(arr.Elements, o)

    for _, el := range node.Elements[1:] {
        o = Eval(el, env)
        if isError(o) { return o }

        if o.Type() != arr.ElementType {
            return createError(
                TypeMismatchError,
                "%s in fixed-type %s of %s",
                o.Type(), object.ArrayType, arr.ElementType)
        }

        arr.Elements = append(arr.Elements, o)
    }

    return arr
}


func isValueObject(obj object.Object) bool {
    _, ok := obj.(object.ValueObject)
    return ok
}

func unwrapReturn(obj object.Object) object.Object {
    if ret, ok := obj.(*object.Return); ok { return ret.Value }
    return obj
}

func createError(errKind string, msg string, args ...any) *object.Error {
    return &object.Error{Message: errKind + ": " + fmt.Sprintf(msg, args...)}
}

func isError(obj object.Object) bool { return obj.Type() == object.ErrorType }
