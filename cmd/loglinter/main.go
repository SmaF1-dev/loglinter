package main

import (
	"github.com/SmaF1-dev/loglinter/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := analyzer.Config{
		SensitiveKeywords: []string{"password", "api_key", "token", "secret", "key", "auth", "env", "environment"},
	}
	singlechecker.Main(analyzer.NewAnalyzer(cfg))
}
