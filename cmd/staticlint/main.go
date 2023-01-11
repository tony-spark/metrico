package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

func main() {
	multichecker.Main(
		printf.Analyzer,
		structtag.Analyzer,
	)
}
