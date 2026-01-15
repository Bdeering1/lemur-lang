package repl

import (
    "bufio"
    "fmt"
    "io"
    "strings"

    "lemur/lexer"
    "lemur/parser"
    "lemur/token"
)

const Prompt = "=> "

const (
    None uint = iota
    Lexer
    Parser
    Stringify
)

func Start(in io.Reader, out io.Writer) {
    fmt.Printf("Welcome to lemur alpha REPL, I'm glad you're here!\n")
    fmt.Printf("Lemur code isn't executable yet, but feel free to explore the lexer and parser!\n")
    fmt.Printf("Please choose a mode:\n")
    fmt.Printf("  'l' for lexer output\n")
    fmt.Printf("  'p' for parser (AST) output\n")
    fmt.Printf("  's' for parsed string output:\n")

    mode := None
    scanner := bufio.NewScanner(in)
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

        if mode == Lexer {
            lex(res)
            continue
        } else if mode == Parser {
            parse(res, false)
        } else {
            parse(res, true)
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
        fmt.Printf("%s\n", program)
    } else {
        fmt.Printf("%s", program.PrintAST())
    }

    errors := p.Errors()
    if len(errors) == 0  { return }

    fmt.Printf("%d parser errors:\n", len(errors))
    for _, msg := range errors {
        fmt.Printf("\t%q\n", msg)
    }
    fmt.Println()
}
