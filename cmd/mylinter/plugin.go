package main

import (
    "github.com/selectel/pkg/analyzer"
    "golang.org/x/tools/go/analysis"
)

func AnalyzerPlugin() []*analysis.Analyzer {
    return []*analysis.Analyzer{analyzer.Analyzer}
}