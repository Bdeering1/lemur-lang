package main

import (
    "os"
    "lemur/api"
)

func main() {
    fi, err := os.Stdin.Stat()
    if err != nil { panic(err) }

    if fi.Mode() & os.ModeNamedPipe != 0 {
        api.EvalFromReader(os.Stdin)
        return
    }

    if len(os.Args) > 1 {
        api.EvalFromFile(os.Args[1])
        return
    }

    api.StartREPL(os.Stdin)
}
