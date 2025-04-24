# Lemur

Lemur is a simple programming language adapted from Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/).

### Syntax features

Declaring variables:
```
let a = 5;
```

Conditionals:
```
if a < b {
  return a
} else {
  return b
}
```

First-class functions:
```
let f = fn(a, b) {
  a + b;
}
```
