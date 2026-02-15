# Lemur

Lemur is an experimental language adapted from Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/).

The language currently supports the following features:
- string, integer, boolean, and array types
- basic logical and arithmentic operations
- variable assignment with implicit typing
- if/else expressions
- first class functions with implicit or explicit returns
- builtin functions for arrays and strings
  - len, first, last, head, tail, push
- interactive REPL with code evaluation + optional lexer and parser output

Syntax sample:
```
let a = 5
let b = a * 2

let max = fn(x, y) { if x > y { x } else { y } }
let c = max(a, b)
```
