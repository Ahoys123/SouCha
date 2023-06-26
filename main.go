package main

import (
	"fmt"
	"io"
	"main/parser"
	"os"
)

func main() {
	l := parser.NewLanguage(readFile("test.soch"))
	// a{b c}d > b
	fmt.Println(l.Evolve("hire this me?"))
}

func readFile(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		return "os.Open error"
	}

	lines, err := io.ReadAll(f)
	if err != nil {
		return "io.ReadAll error"
	}

	return string(lines)
}

//[+vowel] > / _#

// [+stop][S:+stop] > [S:];
// [+stop][N:+nasal] > [N:];

// {p t k(ʷ)} {b d g(ʷ)} > {pʰ tʰ k(ʷ)ʰ} {p t k(ʷ)}
