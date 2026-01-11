package verifier

import (
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/testing/types"
)

type Report struct {
	Timestamp   time.Time            `json:"timestamp"`
	Server      string               `json:"server"`
	ConfigFile  string               `json:"config_file"`
	ProfileName string               `json:"profile_name"`
	Duration    time.Duration        `json:"duration_ms"`
	Summary     Summary              `json:"summary"`
	Checks      []*types.CheckResult `json:"checks"`
	ExitCode    int                  `json:"exit_code"`
}

type Summary struct {
	Total    int `json:"total"`
	Passed   int `json:"passed"`
	Warnings int `json:"warnings"`
	Failed   int `json:"failed"`
}

func NewReport(server string, configFile string, profileName string) *Report {
	return &Report{
		Timestamp:   time.Now(),
		Server:      server,
		ConfigFile:  configFile,
		ProfileName: profileName,
		Checks:      []*types.CheckResult{},
		Summary: Summary{
			Total:    0,
			Passed:   0,
			Warnings: 0,
			Failed:   0,
		},
		ExitCode: 0,
	}
}

func (rep *Report) AddResult(result *types.CheckResult) {
	rep.Checks = append(rep.Checks, result)
	rep.Summary.Total++
	switch result.Status {
	case types.StatusPass:
		rep.Summary.Passed++
	case types.StatusWarning:
		rep.Summary.Warnings++
	case types.StatusFail:
		rep.Summary.Failed++
	}
}

func (rep *Report) ComputeExitCode(warningsOk bool) int {
	if rep.Summary.Failed > 0 {
		return 1
	}
	if rep.Summary.Warnings > 0 && !warningsOk {
		return 1
	}
	return 0
}

func (rep *Report) GetStatus() string {
	if rep.Summary.Failed > 0 {
		return "FAILED"
	}
	if rep.Summary.Warnings > 0 {
		return "PASSED WITH WARNINGS"
	}
	return "PASSED"
}

func (rep *Report) GetStatusEmoji() string {
	if rep.Summary.Failed > 0 {
		return "✗"
	}
	if rep.Summary.Warnings > 0 {
		return "⚠"
	}
	return "✓"
}

func (rep *Report) DurationString() string {
	return fmt.Sprintf("%.3fs", rep.Duration.Seconds())
}

func (rep *Report) SummaryString() string {
	return fmt.Sprintf("%d passed, %d warnings, %d failed",
		rep.Summary.Passed, rep.Summary.Warnings, rep.Summary.Failed)
}
