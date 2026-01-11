package verifier

import (
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/testing/types"
)

type Reporter interface {
	Print(result *types.CheckResult)
	PrintFinal(report *Report)
}

type ConsoleReporter struct {
	mode types.OutputMode
}

func NewConsoleReporter(mode types.OutputMode) *ConsoleReporter {
	return &ConsoleReporter{
		mode: mode,
	}
}

func (r *ConsoleReporter) Print(result *types.CheckResult) {
	if r.mode == types.OutputQuiet {
		return
	}

	icon := "✓"
	status := ""
	switch result.Status {
	case types.StatusFail:
		icon = "✗"
		status = " [FAILED]"
	case types.StatusWarning:
		icon = "⚠"
	}

	duration := ""
	if result.Duration > 0 {
		duration = fmt.Sprintf(" [%6v]", time.Duration(result.Duration).Round(time.Millisecond))
	}

	fmt.Printf("  %s %-30s%s%s\n", icon, result.Check, duration, status)
	if result.Status == types.StatusWarning && result.Message != "" {
		fmt.Printf("    Warning: %s\n", result.Message)
	}
}

func (r *ConsoleReporter) PrintFinal(report *Report) {
	if r.mode == types.OutputQuiet {
		return
	}

	fmt.Printf("\nResults: %s\n", report.SummaryString())
	fmt.Printf("Duration: %s\n", report.DurationString())

	fmt.Printf("\n%s %s\n", report.GetStatusEmoji(), report.GetStatus())

	if r.mode == types.OutputVerbose {
		fmt.Printf("\nDetailed Results:\n")
		for _, check := range report.Checks {
			r.printVerboseCheck(check)
		}
	}
}

func (r *ConsoleReporter) printVerboseCheck(result *types.CheckResult) {
	if result.Status == types.StatusSkip {
		return
	}

	timestamp := result.Timestamp.Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] Running %s check...\n", timestamp, result.Check)

	if result.Status == types.StatusPass {
		fmt.Printf("  ✓ %s\n", result.Message)
	} else if result.Status == types.StatusWarning {
		fmt.Printf("  ⚠ Warning: %s\n", result.Message)
	} else if result.Status == types.StatusFail {
		fmt.Printf("  ✗ Failed: %s\n", result.Message)
		if result.Error != "" {
			fmt.Printf("    Error: %s\n", result.Error)
		}
	}

	if result.Duration > 0 {
		fmt.Printf("  Duration: %v\n", time.Duration(result.Duration).Round(time.Millisecond))
	}

	if len(result.Details) > 0 {
		fmt.Printf("  Details:\n")
		for key, value := range result.Details {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}
}
