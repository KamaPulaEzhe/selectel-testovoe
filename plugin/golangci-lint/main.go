package main

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"golang.org/x/tools/go/analysis"

	"github.com/selectel/pkg/analyzer"
)

func New(conf *config.Config) []*goanalysis.Linter {
	return []*goanalysis.Linter{
		goanalysis.NewLinter(
			analyzer.Analyzer.Name,
			analyzer.Analyzer.Doc,
			[]*analysis.Analyzer{analyzer.Analyzer},
			nil,
		).WithContextSetter(func(lintCtx *linter.Context) {}),
	}
}
