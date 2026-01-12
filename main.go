package main

import (
	"github.com/YuitoSato/gounion/gounion"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(gounion.Analyzer)
}
