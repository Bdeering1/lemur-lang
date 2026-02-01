package repl

import (
    "bufio"
    "fmt"
    "io"
    "strings"

    "lemur/eval"
    "lemur/lexer"
    "lemur/object"
    "lemur/parser"
    "lemur/token"
)

const Prompt = "=> "

const (
    None uint = iota
    Lexer
    Parser
    Stringify
    Evaluate
)

func Start(in io.Reader, out io.Writer) {
    fmt.Printf("Welcome to lemur alpha REPL, I'm glad you're here!\n")
    fmt.Printf("Please choose a mode:\n")
    fmt.Printf("  'l' for lexer output\n")
    fmt.Printf("  'p' for parser (AST) output\n")
    fmt.Printf("  's' for parsed string output\n")
    fmt.Printf("  'e' for code evaluation (default)\n")

    mode := None
    scanner := bufio.NewScanner(in)
    env := object.CreateEnvironment()

    for {
        res := prompt(scanner)
        if res == "" || res == "q" || res == "quit" { break }

        if res == "l" || res == "lexer" {
            fmt.Printf("<lexer mode>\n")
            mode = Lexer
            continue
        }
        if res == "p" || res == "parser" {
            fmt.Printf("<parser mode>\n")
            mode = Parser
            continue
        }
        if res == "s" || res == "string" {
            fmt.Printf("<string mode>\n")
            mode = Stringify
            continue
        }
        if res == "e" || res == "eval" {
            fmt.Printf("<eval mode>\n")
            mode = Evaluate
            continue
        }

        if mode == Lexer {
            lex(res)
            continue
        } else if mode == Parser {
            parse(res, false)
        } else {
            evaluate(res, env)
        }
    }

}

func prompt(scanner *bufio.Scanner) string {
    fmt.Print(Prompt)
    if !scanner.Scan() { return "" }

    res := strings.TrimSpace(scanner.Text())
    return res
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

func evaluate(input string, env *object.Environment) {
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

func printParserErrors(errors []string) {
    fmt.Printf("Failed to parse (%d errors):\n", len(errors))
    for _, msg := range errors {
        fmt.Printf("  Error: %s\n", msg)
    }
}
