package eval

import "lemur/object"

const (
    Len   = "len"
    First = "first"
    Last  = "last"
    Head  = "head"
    Tail  = "tail"
    Push  = "push"
)

var builtins = map[string]object.Builtin{
    Len: func(args ...object.Object) object.Object {
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
    First: func(args ...object.Object) object.Object {
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
    Last: func(args ...object.Object) object.Object {
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
    Head: func(args ...object.Object) object.Object {
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
    Tail: func(args ...object.Object) object.Object {
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
    Push: func(args ...object.Object) object.Object {
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
}
