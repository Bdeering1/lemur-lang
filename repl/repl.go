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

func Start(in io.Reader, out io.Writer) {
    fmt.Printf("Welcome to lemur alpha REPL, I'm glad you're here!\n")
    fmt.Printf("Lemur code isn't executable yet, but feel free to explore the lexer and parser!\n")
    fmt.Printf("Please type 'l' for lexer output and 'p' for parser output:\n")

    mode := "none"
    scanner := bufio.NewScanner(in)
    for {
        res := prompt(scanner)
        if res == "" || res == "q" || res == "quit" { break }

        if res == "l" || res == "lexer" {
            fmt.Printf("<lexer mode>\n")
            mode = "lex"
            continue
        }
        if res == "p" || res == "parser" {
            fmt.Printf("<parser mode>\n")
            mode = "parse"
            continue
        }

        if mode == "lex" {
            lex(res)
            continue
        }
        parse(res)
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

func parse(input string) {
    input = input + "\x00"
    l := lexer.New(input)
    p := parser.New(l)

    program := p.ParseProgram()
    fmt.Printf("%s", program.PrintAST())

    errors := p.Errors()
    if len(errors) == 0  { return }

    fmt.Printf("%d parser errors:", len(errors))
    for _, msg := range errors {
        fmt.Printf("\t%q", msg)
    }
    fmt.Println()
}
