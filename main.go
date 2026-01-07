package main

import (
    "os"
    "lemur/repl"
)

func main() {
    repl.Start(os.Stdin, os.Stdout)
}
