package api

import (
    "fmt"
    "io"
    "os"

    "lemur/eval"
    "lemur/lexer"
    "lemur/parser"
    "lemur/object"
    "lemur/token"
)

func EvalFromReader(in io.Reader) {
    b, err := io.ReadAll(in)
    if err != nil || len(b) == 0 { return }

    input := string(b)
    env := object.CreateEnvironment()
    runEval(input, env)
}

func EvalFromFile(fname string) {
    b, err := os.ReadFile(fname)
    if err != nil || len(b) == 0 { return }

    input := string(b)
    env := object.CreateEnvironment()
    runEval(input, env)
}

func runEval(input string, env *object.Environment) {
    input = input + "\x00"
    l := lexer.New(input)
    p := parser.New(l)

    program := p.ParseProgram()
    if len(p.Errors()) != 0  {
        printParserErrors(p.Errors())
        return
    }

    evaluated := eval.Eval(program, env)
    fmt.Println(evaluated.String())
}

func lex(input string) {
    l := lexer.New(input)
    for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
        fmt.Printf("%+v\n", tok)
    }
}

func parse(input string, stringify bool) {
    input = input + "\x00"
    l := lexer.New(input)
    p := parser.New(l)

    program := p.ParseProgram()
    if stringify {
        fmt.Printf("%s\n", program.String())
    } else {
        fmt.Printf("%s", program.PrintAST())
    }

    if len(p.Errors()) == 0  { return }
    printParserErrors(p.Errors())
}

func printParserErrors(errors []string) {
    fmt.Printf("Failed to parse (%d errors):\n", len(errors))
    for _, msg := range errors {
        fmt.Printf("  Error: %s\n", msg)
    }
}
