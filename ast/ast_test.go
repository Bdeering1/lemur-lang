package ast

import (
    "testing"

    "lemur/token"
)

func TestString(t *testing.T) {
    program := &Program{
        &LetStatement{
            Token: token.New(token.Let, "let"),
            Name: &Identifier{
                Token: token.New(token.Ident, "myVar"),
                Value: "myVar",
            },
            Value: &Identifier{
                Token: token.New(token.Ident, "anotherVar"),
                Value: "anotherVar",
            },
        },
    }

    if program.String() != "let myVar = anotherVar;" {
        t.Errorf("program.String() is incorrect (got %q)", program.String())
    }
}
