# Lemur

Lemur is an interpreted language in early development, adapted from Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/). It currently consists of a handwritten lexer and parser.

The language currently supports the following features:
- integer, boolean, and string literals
- basic arithmentic operations
- variable assignment with implicit typing
- if/else expressions
- first class functions

Syntax sample:
```
let a = 5
let b = a * 2

let max = fn(x, y) { if x > y { x } else { y } }
let c = max(a, b)
```
