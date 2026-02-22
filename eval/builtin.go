package eval

import (
    "fmt"

    "lemur/object"
)

const (
    Puts    = "puts"
    Len     = "len"
    First   = "first"
    Last    = "last"
    Head    = "head"
    Tail    = "tail"
    Push    = "push"
    Iter    = "iter"
    Collect = "collect"
)

var builtins = map[string]object.Builtin{
    Puts: func(args []object.Object) object.Object {
        for _, arg := range args {
            fmt.Printf("%s ", arg)
        }
        fmt.Println()

        return Null
    },
    Len: func(args []object.Object) object.Object {
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", Len)
        }

        switch input := args[0].(type) {
        case *object.Array:
            return object.Integer(len(input.Elements))
        case object.String:
            return object.Integer(len(input))
        default:
            return createError(ArgumentTypesError, "%s(%s)", Len,  input.Type())
        }
    },
    First: func(args []object.Object) object.Object {
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", First)
        }

        switch input := args[0].(type) {
        case *object.Array:
            if len(input.Elements) == 0 { return Null }
            return input.Elements[0]
        case object.String:
            if len(input) == 0 { return Null }
            return object.String(input[0])
        default:
            return createError(ArgumentTypesError, "%s(%s)", First,  input.Type())
        }
    },
    Last: func(args []object.Object) object.Object {
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", Last)
        }

        switch input := args[0].(type) {
        case *object.Array:
            if len(input.Elements) == 0 { return Null }
            return input.Elements[len(input.Elements) - 1]
        case object.String:
            if len(input) == 0 { return Null }
            return object.String(input[len(input) - 1])
        default:
            return createError(ArgumentTypesError, "%s(%s)", Last,  input.Type())
        }
    },
    Head: func(args []object.Object) object.Object {
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", Head)
        }

        switch input := args[0].(type) {
        case *object.Array:
            if len(input.Elements) < 2 { return &object.Array{Elements: []object.Object{}} }
            return &object.Array{Elements: input.Elements[0:len(input.Elements) - 1]}
        case object.String:
            if len(input) < 2 { return object.String("") }
            return object.String(input[:len(input) - 1])
        default:
            return createError(ArgumentTypesError, "%s(%s)", Head,  input.Type())
        }
    },
    Tail: func(args []object.Object) object.Object {
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", Tail)
        }

        switch input := args[0].(type) {
        case *object.Array:
            if len(input.Elements) < 2 { return &object.Array{Elements: []object.Object{}} }
            return &object.Array{Elements: input.Elements[1:len(input.Elements)]}
        case object.String:
            if len(input) < 2 { return object.String("") }
            return object.String(input[1:])
        default:
            return createError(ArgumentTypesError, "%s(%s)", Tail, input.Type())
        }
    },
    Push: func(args []object.Object) object.Object {
        if len(args) != 2 {
            return createError(ArgumentMistmatchError, "%s", Push)
        }

        switch col := args[0].(type) {
        case *object.Array:
            if col.ElementType == object.NullType {
                col.ElementType = args[1].Type()
            } else if args[1].Type() != col.ElementType {
                return createError(
                    TypeMismatchError,
                    "%s(Array[%v], %v)",
                    Push, col.ElementType, args[1].Type())
            }

            arr := append(col.Elements, args[1])
            return &object.Array{Elements: arr}
        case object.String:
            obj, ok := args[1].(object.String)
            if !ok {
                return createError(
                    ArgumentTypesError,
                    "%s(%s, %v)",
                    Push, object.StringType, obj.Type())
            }

            return object.String(col+ obj)
        default:
            return createError(
                ArgumentTypesError,
                "%s(%v, %v)",
                Push, col.Type(), args[1].Type())
        }
    },
    Iter: func(args []object.Object) object.Object {
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", Iter)
        }

        switch col := args[0].(type) {
        case *object.Array:
            idx := 0

            var f object.Builtin
            f = func(args []object.Object) object.Object {
                if idx >= len(col.Elements) {
                    return object.Tuple{
                        zero(col.ElementType),
                        object.Boolean(false),
                    }
                }

                res := col.Elements[idx]; idx++
                return object.Tuple{
                    res,
                    object.Boolean(true),
                }
            }
            return f

        case object.String:
            idx := 0

            var f object.Builtin
            f = func(args []object.Object) object.Object {
                if idx >= len(col) {
                    return object.Tuple{
                        zero(object.StringType),
                        object.Boolean(false),
                    }
                }

                res := object.String(string(col[idx])); idx++
                return object.Tuple{
                    res,
                    object.Boolean(true),
                }
            }
            return f

        default:
            return createError(
                ArgumentTypesError,
                "%s(%v)",
                Iter, col.Type())
        }
    },
    Collect: func(args []object.Object) object.Object { // only creates arrays for now
        if len(args) != 1 {
            return createError(ArgumentMistmatchError, "%s", Collect)
        }

        next, ok := args[0].(object.Builtin)
        if !ok {
            return createError(ArgumentTypesError, "%s(%s)", Collect, args[0].Type())
        }

        arr := &object.Array{
            Elements: []object.Object{},
            ElementType: object.NullType,
        }

        t := next([]object.Object{}).(object.Tuple)
        val := t[0]
        ok = bool(t[1].(object.Boolean))

        if !ok { return arr }
        arr.ElementType = val.Type()
        arr.Elements = append(arr.Elements, val)

        for {
            t = next([]object.Object{}).(object.Tuple)
            val = t[0]
            ok = bool(t[1].(object.Boolean))

            if !ok { return arr }
            arr.Elements = append(arr.Elements, val)
        }
    },
}

func zero(t object.ObjectType) object.Object {
    switch t {
    case object.StringType:
        return object.String("")
    case object.IntegerType:
        return object.Integer(0)
    case object.BooleanType:
        return object.Boolean(false)
    }

    return Null
}
