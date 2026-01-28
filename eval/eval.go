package eval

import (
    "lemur/ast"
    "lemur/object"
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

    case *ast.ExpressionStatement:
        return Eval(node.Value)

    case *ast.ConditionalExpression:
        return evalConditionalExpression(node)

    case *ast.InfixExpression:
        left := Eval(node.Left)
        right := Eval(node.Right)
        return evalInfixExpression(node.Operator, left, right)

    case *ast.PrefixExpression:
        right := Eval(node.Right)
        return evalPrefixOperator(node.Operator, right)
        
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.BooleanLiteral:
        return createBooleanObject(node.Value)
    }

    return nil
}

func evalBlock(block []ast.Statement) object.Object {
    var res object.Object

    for _, stmt := range block {
        res = Eval(stmt)
    }

    return res
}

func evalConditionalExpression(ce *ast.ConditionalExpression) object.Object {
    cond := Eval(ce.Condition)

    if cond == True { return Eval(ce.Consequence) }
    if cond == False && ce.Alternative != nil { return Eval(ce.Alternative) }

    return Null
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
    if left.Type() != right.Type() { return Null } // raise error

    switch {
    case left.Type() == object.IntegerType:
        return evalIntegerInfixExpression(operator, left, right)
    case left.Type() == object.BooleanType:
        return evalBooleanInfixExpression(operator, left, right)
    default:
        return Null // raise error
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
        return Null // raise error
    }
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
    switch operator {
    case "==":
        return createBooleanObject(left == right)
    case "!=":
        return createBooleanObject(left != right)
    default:
        return Null
    }
}

func evalPrefixOperator(operator string, right object.Object) object.Object {
    switch operator {
    case "!":
        return evalBangPrefix(right)
    case "-":
        return evalMinusPrefix(right)        
    default:
        return Null // raise error
    }
}

func evalBangPrefix(right object.Object) object.Object {
    switch right {
    case True:
        return False
    case False:
        return True
    default:
        return Null // raise error
    }
}

func evalMinusPrefix(right object.Object) object.Object {
    if right.Type() != object.IntegerType { return Null } // raise error
    
    val := right.(*object.Integer).Value
    return &object.Integer{Value: -val}
}

func createBooleanObject(val bool) object.Object{
    if val { return True } else { return False }
}
