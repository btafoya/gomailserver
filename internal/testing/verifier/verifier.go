package verifier

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/testing/checks"
	"github.com/btafoya/gomailserver/internal/testing/types"
)

type Verifier struct {
	config       *Config
	serverConfig *types.ServerConfig
	report       *Report
	registry     *checks.Registry
	rateLimiter  *time.Ticker
}

func NewVerifier(cfg *Config, serverConfig *types.ServerConfig) *Verifier {
	return &Verifier{
		config:       cfg,
		serverConfig: serverConfig,
		report:       NewReport(serverConfig.SMTPHost, cfg.ConfigFile, cfg.ProfileName),
		registry:     checks.NewRegistry(),
	}
}

func (v *Verifier) InitializeChecks() {
	v.registry.RegisterAll()
}

func (v *Verifier) Run(ctx context.Context) (*Report, error) {
	startTime := time.Now()

	fmt.Printf("Running verification for %s...\n", v.serverConfig.SMTPHost)

	v.InitializeChecks()

	allChecks := v.registry.GetAll()
	for _, check := range allChecks {
		result, err := check.Run(ctx, v.serverConfig)
		if err != nil {
			fmt.Printf("Error running %s: %v\n", check.Name(), err)
		} else if result != nil {
			v.report.AddResult(result)
		}
	}

	v.report.Duration = time.Since(startTime)
	v.report.ExitCode = v.report.ComputeExitCode(v.config.WarningsOk)

	return v.report, nil
}

func (v *Verifier) PrintFinalSummary() {
	if v.config.OutputMode == types.OutputQuiet {
		return
	}

	fmt.Printf("\nResults: %s\n", v.report.SummaryString())
	fmt.Printf("Duration: %s\n", v.report.DurationString())

	fmt.Printf("\n%s %s\n", v.report.GetStatusEmoji(), v.report.GetStatus())
}

func (v *Verifier) GetReport() *Report {
	return v.report
}
