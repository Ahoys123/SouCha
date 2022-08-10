package main

import (
	"main/parser"
)

func main() {
	parser.NewLanguage("a(b) (b)c > d", "a bc")
	// a{b c}d > b
}
