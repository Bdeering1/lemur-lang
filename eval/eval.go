package eval

import (
    "lemur/ast"
    "lemur/object"
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {

    case ast.Program:
        return evalBlock(node)

    case *ast.BlockStatement:
        return evalBlock(node.Statements)

    case *ast.ExpressionStatement:
        return Eval(node.Value)

    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}
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
