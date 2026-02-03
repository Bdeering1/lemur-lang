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

    api.StartREPL(os.Stdin)
}
