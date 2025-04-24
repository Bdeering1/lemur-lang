package main

import (
    "fmt"
    "os"
    "lemur/repl"
)

func main() {
    fmt.Printf("Welcome to lemur alpha, we're glad you're here!\n")
    fmt.Printf("Feel free to start by entering some commands\n")
    repl.Start(os.Stdin, os.Stdout)
}
