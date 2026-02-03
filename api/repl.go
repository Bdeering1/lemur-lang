package api

import (
    "bufio"
    "fmt"
    "io"
    "strings"

    "lemur/object"
)

const Prompt = "=> "

const (
    None uint = iota
    Lexer
    Parser
    Stringify
    Evaluate
)

func StartREPL(in io.Reader) {
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
            runEval(res, env)
        }
    }

}

func prompt(scanner *bufio.Scanner) string {
    fmt.Print(Prompt)
    if !scanner.Scan() { return "" }

    res := strings.TrimSpace(scanner.Text())
    return res
}
