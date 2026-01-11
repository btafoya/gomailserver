package verifier

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/btafoya/gomailserver/internal/testing/types"
)

type JSONReporter struct {
	filePath string
}

func NewJSONReporter(filePath string) *JSONReporter {
	return &JSONReporter{
		filePath: filePath,
	}
}

func (r *JSONReporter) Print(result *types.CheckResult) {
	return
}

func (r *JSONReporter) PrintFinal(report *Report) {
	if r.filePath == "" {
		return
	}

	output := struct {
		Timestamp   time.Time            `json:"timestamp"`
		Server      string               `json:"server"`
		ConfigFile  string               `json:"config_file"`
		ProfileName string               `json:"profile_name"`
		Duration    time.Duration        `json:"duration_ms"`
		Summary     Summary              `json:"summary"`
		Checks      []*types.CheckResult `json:"checks"`
		ExitCode    int                  `json:"exit_code"`
		Status      string               `json:"status"`
	}{
		Timestamp:   report.Timestamp,
		Server:      report.Server,
		ConfigFile:  report.ConfigFile,
		ProfileName: report.ProfileName,
		Duration:    report.Duration,
		Summary:     report.Summary,
		Checks:      report.Checks,
		ExitCode:    report.ExitCode,
		Status:      report.GetStatus(),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON report: %v\n", err)
		return
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON report: %v\n", err)
		return
	}

	fmt.Printf("JSON report saved to: %s\n", r.filePath)
}

func (r *JSONReporter) PrintToWriter(w io.Writer, report *Report) error {
	output := struct {
		Timestamp   time.Time            `json:"timestamp"`
		Server      string               `json:"server"`
		ConfigFile  string               `json:"config_file"`
		ProfileName string               `json:"profile_name"`
		Duration    time.Duration        `json:"duration_ms"`
		Summary     Summary              `json:"summary"`
		Checks      []*types.CheckResult `json:"checks"`
		ExitCode    int                  `json:"exit_code"`
		Status      string               `json:"status"`
	}{
		Timestamp:   report.Timestamp,
		Server:      report.Server,
		ConfigFile:  report.ConfigFile,
		ProfileName: report.ProfileName,
		Duration:    report.Duration,
		Summary:     report.Summary,
		Checks:      report.Checks,
		ExitCode:    report.ExitCode,
		Status:      report.GetStatus(),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

type HTMLReporter struct {
	filePath string
}

func NewHTMLReporter(filePath string) *HTMLReporter {
	return &HTMLReporter{
		filePath: filePath,
	}
}

func (r *HTMLReporter) Print(result *types.CheckResult) {
	return
}

func (r *HTMLReporter) PrintFinal(report *Report) {
	if r.filePath == "" {
		return
	}

	html := r.generateHTML(report)

	if err := os.WriteFile(r.filePath, []byte(html), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing HTML report: %v\n", err)
		return
	}

	fmt.Printf("HTML report saved to: %s\n", r.filePath)
}

func (r *HTMLReporter) generateHTML(report *Report) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoMailServer Verification Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            border-radius: 8px;
            margin-bottom: 30px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header h1 {
            margin: 0 0 10px 0;
            font-size: 28px;
        }
        .header .status {
            font-size: 18px;
            font-weight: bold;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .summary-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        .summary-card h3 {
            margin: 0 0 10px 0;
            color: #666;
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .summary-card .value {
            font-size: 32px;
            font-weight: bold;
            color: #333;
        }
        .summary-card .value.pass { color: #10b981; }
        .summary-card .value.warning { color: #f59e0b; }
        .summary-card .value.fail { color: #ef4444; }
        .section {
            background: white;
            padding: 25px;
            border-radius: 8px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        .section h2 {
            margin: 0 0 20px 0;
            color: #333;
            font-size: 20px;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .check-item {
            padding: 15px;
            margin-bottom: 15px;
            border-radius: 6px;
            border-left: 4px solid #ccc;
        }
        .check-item.pass {
            background: #ecfdf5;
            border-left-color: #10b981;
        }
        .check-item.warning {
            background: #fffbeb;
            border-left-color: #f59e0b;
        }
        .check-item.fail {
            background: #fef2f2;
            border-left-color: #ef4444;
        }
        .check-item.skip {
            background: #f9fafb;
            border-left-color: #6b7280;
        }
        .check-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        .check-name {
            font-weight: bold;
            font-size: 16px;
        }
        .check-status {
            padding: 4px 12px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: bold;
            text-transform: uppercase;
        }
        .check-status.pass {
            background: #10b981;
            color: white;
        }
        .check-status.warning {
            background: #f59e0b;
            color: white;
        }
        .check-status.fail {
            background: #ef4444;
            color: white;
        }
        .check-status.skip {
            background: #6b7280;
            color: white;
        }
        .check-message {
            margin: 10px 0;
            color: #555;
        }
        .check-details {
            margin-top: 10px;
            padding: 10px;
            background: #f9fafb;
            border-radius: 4px;
            font-size: 13px;
            color: #666;
        }
        .check-details dt {
            font-weight: bold;
            margin-top: 8px;
        }
        .check-details dd {
            margin-left: 0;
            margin-bottom: 8px;
        }
        .timestamp {
            text-align: center;
            color: #666;
            font-size: 12px;
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>GoMailServer Verification Report</h1>
        <div class="status">` + report.GetStatus() + `</div>
    </div>

    <div class="summary">
        <div class="summary-card">
            <h3>Total Checks</h3>
            <div class="value">` + fmt.Sprintf("%d", report.Summary.Total) + `</div>
        </div>
        <div class="summary-card">
            <h3>Passed</h3>
            <div class="value pass">` + fmt.Sprintf("%d", report.Summary.Passed) + `</div>
        </div>
        <div class="summary-card">
            <h3>Warnings</h3>
            <div class="value warning">` + fmt.Sprintf("%d", report.Summary.Warnings) + `</div>
        </div>
        <div class="summary-card">
            <h3>Failed</h3>
            <div class="value fail">` + fmt.Sprintf("%d", report.Summary.Failed) + `</div>
        </div>
    </div>

    <div class="section">
        <h2>Server Information</h2>
        <dl class="check-details">
            <dt>Server</dt>
            <dd>` + report.Server + `</dd>
            <dt>Config File</dt>
            <dd>` + report.ConfigFile + `</dd>
            <dt>Profile</dt>
            <dd>` + report.ProfileName + `</dd>
            <dt>Duration</dt>
            <dd>` + report.DurationString() + `</dd>
            <dt>Exit Code</dt>
            <dd>` + fmt.Sprintf("%d", report.ExitCode) + `</dd>
        </dl>
    </div>

    <div class="section">
        <h2>Check Results</h2>
        ` + r.generateChecksHTML(report) + `
    </div>

    <div class="timestamp">
        Generated on ` + report.Timestamp.Format("2006-01-02 15:04:05 MST") + `
    </div>
</body>
</html>`
}

func (r *HTMLReporter) generateChecksHTML(report *Report) string {
	html := ""
	for _, check := range report.Checks {
		html += r.generateCheckHTML(check)
	}
	return html
}

func (r *HTMLReporter) generateCheckHTML(result *types.CheckResult) string {
	statusClass := string(result.Status)
	if statusClass == "warning" || statusClass == "fail" || statusClass == "pass" || statusClass == "skip" {
		statusClass = string(result.Status)
	}

	return `
        <div class="check-item ` + statusClass + `">
            <div class="check-header">
                <div class="check-name">` + result.Check + `</div>
                <div class="check-status ` + statusClass + `">` + string(result.Status) + `</div>
            </div>
            <div class="check-message">` + result.Message + `</div>
            ` + r.generateDetailsHTML(result) + `
        </div>`
}

func (r *HTMLReporter) generateDetailsHTML(result *types.CheckResult) string {
	if len(result.Details) == 0 {
		return ""
	}

	html := `<div class="check-details">`
	for key, value := range result.Details {
		html += fmt.Sprintf("<dt>%s</dt><dd>%v</dd>", key, value)
	}
	html += `</div>`
	return html
}
