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

    case *ast.PrefixExpression:
        right := Eval(node.Right)
        return evalPrefixOperator(node.Operator, right)
        
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.BooleanLiteral:
        if node.Value { return True } else { return False }
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

func evalPrefixOperator(operator string, right object.Object) object.Object {
    switch operator {
    case "!":
        return evalBangPrefix(right)
    case "-":
        return evalMinusPrefix(right)        
    default:
        return Null
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
