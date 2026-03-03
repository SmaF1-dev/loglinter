package analyzer_test

import (
	"testing"

	"github.com/SmaF1-dev/loglinter/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.NewAnalyzer(analyzer.Config{}), "loglint")
}
