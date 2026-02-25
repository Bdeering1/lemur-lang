# Lemur

Lemur is an experimental language inspired by Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/).

The language currently supports the following features:
- string, integer, boolean, array, and hash map types
- basic logical and arithmentic operations
- assignment with implicit typing, multiple assignment and tuple expressions
- if/else expressions
- first class functions with implicit or explicit returns
- pipeline expressions
- builtin functions
  - puts, len, first, last, head, tail, push, iter, map, collect
- EOL comments
- interactive REPL with code evaluation + optional lexer and parser output

Syntax sample:
```ruby
let double = fn(col) {
    iter(col), fn(x){ x * 2 }
    |> map
    |> collect
}

let arr = [1, 2, 3]
puts(double(arr)) # [2, 4, 6]
```

## Usage

```sh
git clone https://github.com/Bdeering1/lemur-lang.git
cd lemur-lang && go build

./lemur # REPL
./lemur my_file.txt
```
