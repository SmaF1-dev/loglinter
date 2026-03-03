package main

import (
	"github.com/SmaF1-dev/loglinter/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(conf interface{}) ([]*analysis.Analyzer, error) {
	var cfg analyzer.Config

	if setting, ok := conf.(map[string]interface{}); ok {
		if keywords, ok := setting["sensitive_keywords"].([]interface{}); ok {
			for _, kw := range keywords {
				if s, ok := kw.(string); ok {
					cfg.SensitiveKeywords = append(cfg.SensitiveKeywords, s)
				}
			}
		}
	}

	return []*analysis.Analyzer{analyzer.NewAnalyzer(cfg)}, nil
}
