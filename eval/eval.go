package eval

import (
    "fmt"

    "lemur/ast"
    "lemur/object"
)

const (
    InfixNotImplementedError = "no infixes implemented for type"
    InvalidConditionError = "invalid condition"
    TypeMismatchError = "type mismatch"
    UnknownOperatorError = "unknown operator"
    InternalErrorPostfix = " (internal)"
)

var (
    True = &object.Boolean{Value: true}
    False = &object.Boolean{Value: false}
    Null = &object.Null{}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {

    case ast.Program:
        return evalBlock(node)

    case *ast.BlockStatement:
        return evalBlock(node.Statements)

    case *ast.ReturnStatement:
        obj := Eval(node.Value)
        if isError(obj) { return obj }

        return &object.Return{Value: obj}

    case *ast.ExpressionStatement:
        return Eval(node.Value)

    case *ast.ConditionalExpression:
        return evalConditionalExpression(node)

    case *ast.InfixExpression:
        left := Eval(node.Left)
        if isError(left) { return left }

        right := Eval(node.Right)
        if isError(right) { return right }

        return evalInfixExpression(node.Operator, left, right)

    case *ast.PrefixExpression:
        right := Eval(node.Right)
        if isError(right) { return right }

        return evalPrefixOperator(node.Operator, right)
        
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.BooleanLiteral:
        return createBooleanObject(node.Value)
    }

    return nil
}

func evalBlock(block []ast.Statement) object.Object {
    var obj object.Object

    for _, stmt := range block {
        obj = Eval(stmt)

        if obj.Type() == object.ErrorType || obj.Type() == object.ReturnType { return obj }
    }

    return obj
}

func evalConditionalExpression(ce *ast.ConditionalExpression) object.Object {
    cond := Eval(ce.Condition)
    if isError(cond) { return cond }

    if cond == True { return Eval(ce.Consequence) }
    if cond == False {
        if ce.Alternative == nil { return Null } // or default value for type
        return Eval(ce.Alternative)
    }

    return createError(InvalidConditionError, "%s", ce.Condition)
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
    if left.Type() != right.Type() {
        return createError(TypeMismatchError, "%s %s %s", left.Type(), operator, right.Type())
    }

    switch {
    case left.Type() == object.IntegerType:
        return evalIntegerInfixExpression(operator, left, right)
    case left.Type() == object.BooleanType:
        return evalBooleanInfixExpression(operator, left, right)
    default:
        return createError(InfixNotImplementedError, "%s", left.Type())
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
    switch operator {
    case "==":
        return createBooleanObject(left == right)
    case "!=":
        return createBooleanObject(left != right)
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

func createBooleanObject(val bool) object.Object{
    if val { return True } else { return False }
}

func createError(errKind string, msg string, args ...any) *object.Error {
    return &object.Error{Message: errKind + ": " + fmt.Sprintf(msg, args...)}
}

func isError(obj object.Object) bool { return obj.Type() == object.ErrorType }
