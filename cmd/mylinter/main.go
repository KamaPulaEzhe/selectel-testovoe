//go:build !plugin
// +build !plugin

package main

import (
	"github.com/selectel/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
