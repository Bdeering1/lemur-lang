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
```rust
let map = fn(col, f) {
    let iter = fn(col, res) {
        if len(col) == 0 {
            return res
        }
        iter(tail(col), push(res, f(first(col)))
    }

    iter(col, [])
}

let arr = [1, 2, 3]
map(arr, fn(x){ x * 2 }) // [2, 4, 6]
```

## Usage

```sh
git clone https://github.com/Bdeering1/lemur-lang.git
cd lemur-lang && go build

lemur # REPL
lemur my_file.txt
```
