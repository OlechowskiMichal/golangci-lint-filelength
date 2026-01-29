package filelength

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFileOverLimit(t *testing.T) {
	dir := testdataDir(t)
	a := NewAnalyzer(Settings{MaxLines: 300})
	analysistest.Run(t, dir, a, "example")
}

func TestExcludePatterns(t *testing.T) {
	dir := testdataDir(t)
	a := NewAnalyzer(Settings{MaxLines: 300, ExcludePatterns: []string{"mock_*.go"}})
	analysistest.Run(t, dir, a, "excludepattern")
}

func TestDefaultSettings(t *testing.T) {
	a := NewAnalyzerWithDefaults()
	if a.Name != "filelength" {
		t.Errorf("expected analyzer name 'filelength', got %q", a.Name)
	}
}

func TestShouldSkip(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		settings Settings
		want     bool
	}{
		{
			name:     "regular file not skipped",
			filename: "/path/to/main.go",
			settings: Settings{MaxLines: 300},
			want:     false,
		},
		{
			name:     "test file skipped when ExcludeTests is true",
			filename: "/path/to/main_test.go",
			settings: Settings{MaxLines: 300, ExcludeTests: true},
			want:     true,
		},
		{
			name:     "test file not skipped when ExcludeTests is false",
			filename: "/path/to/main_test.go",
			settings: Settings{MaxLines: 300, ExcludeTests: false},
			want:     false,
		},
		{
			name:     "file matching exclude pattern skipped",
			filename: "/path/to/mock_service.go",
			settings: Settings{MaxLines: 300, ExcludePatterns: []string{"mock_*.go"}},
			want:     true,
		},
		{
			name:     "file not matching exclude pattern not skipped",
			filename: "/path/to/service.go",
			settings: Settings{MaxLines: 300, ExcludePatterns: []string{"mock_*.go"}},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkip(tt.filename, tt.settings)
			if got != tt.want {
				t.Errorf("shouldSkip(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func testdataDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
