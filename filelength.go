package filelength

import (
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer(settings Settings) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "filelength",
		Doc:  "Checks that files do not exceed a maximum number of lines.",
		Run:  runWithSettings(settings),
	}
}

func runWithSettings(settings Settings) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		for _, file := range pass.Files {
			filename := pass.Fset.Position(file.Pos()).Filename

			if shouldSkip(filename, settings) {
				continue
			}

			lineCount := pass.Fset.Position(file.End()).Line
			if lineCount > settings.MaxLines {
				pass.Reportf(file.Pos(), "file is %d lines long, which exceeds the maximum of %d lines", lineCount, settings.MaxLines)
			}
		}

		return nil, nil
	}
}

func shouldSkip(filename string, settings Settings) bool {
	base := filepath.Base(filename)

	if settings.ExcludeTests && strings.HasSuffix(base, "_test.go") {
		return true
	}

	for _, pattern := range settings.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			return true
		}
	}

	return false
}

// NewAnalyzerWithDefaults creates an analyzer with default settings for use with analysistest.
func NewAnalyzerWithDefaults() *analysis.Analyzer {
	return NewAnalyzer(Settings{MaxLines: 300})
}
