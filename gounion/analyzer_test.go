package gounion_test

import (
	"testing"

	"github.com/YuitoSato/gounion/gounion"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	// Run tests on all test packages
	// The order matters: union must be analyzed before consumer
	analysistest.Run(t, testdata, gounion.Analyzer,
		"union",
		"consumer",
	)
}
