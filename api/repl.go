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
    fmt.Printf("Welcome to the lemur alpha REPL!\n")
    fmt.Printf("Start typing for code evaluation, or choose another mode:\n")
    fmt.Printf("  l: lexer output\n")
    fmt.Printf("  p: parser output (AST)\n")
    fmt.Printf("  s: parser output (stringified)\n")
    fmt.Printf("  e: code evaluation (default)\n\n")

    mode := None
    scanner := bufio.NewScanner(in)
    env := object.CreateEnvironment()

    for {
        res := prompt(scanner)
        if res == "q" || res == "quit" { break }
        if res == "" { continue }

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
        } else if mode == Stringify {
            parse(res, true)
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
