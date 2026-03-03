package main

import (
	"github.com/SmaF1-dev/loglinter/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(conf interface{}) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.NewAnalyzer()}, nil
}
