package filelength

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("filelength", New)
}

type Settings struct {
	MaxLines        int      `json:"maxLines"`
	ExcludeTests    bool     `json:"excludeTests"`
	ExcludePatterns []string `json:"excludePatterns"`
}

type Plugin struct {
	settings Settings
}

func New(conf any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](conf)
	if err != nil {
		return nil, err
	}

	if s.MaxLines == 0 {
		s.MaxLines = 300
	}

	return &Plugin{settings: s}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		NewAnalyzer(p.settings),
	}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
