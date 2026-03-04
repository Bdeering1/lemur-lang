package lexer

import "testing"

func resetLexer(l *Lexer) { // reset lexer with minimum overhead
    l.pos = 0
    l.nextPos = 1
    l.ch = l.input[0]
}

func BenchmarkOperator(b *testing.B) {
    input := "="
    l := &Lexer{input: input}

    for b.Loop() {
        resetLexer(l)
        l.NextToken()
    }
}

func BenchmarkCompositeOperator(b *testing.B) {
    input := "=="
    l := &Lexer{input: input}

    for b.Loop() {
        resetLexer(l)
        l.NextToken()
    }
}
