package repl

import (
    "bufio"
    "fmt"
    "io"
    "strings"

    "lemur/lexer"
    "lemur/token"
)

const prompt = "=> "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        fmt.Fprint(out, prompt)
        if !scanner.Scan() { return }

        line := strings.TrimSpace(scanner.Text())
        if line == "q" || line == "quit" { break }
        l := lexer.New(line)

        for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
            fmt.Fprintf(out, "%+v\n", tok)
        }
    }
}
