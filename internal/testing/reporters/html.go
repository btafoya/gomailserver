package reporters

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/testing"
)

// HTMLReporter generates HTML reports for test results
type HTMLReporter struct {
	template *template.Template
}

// NewHTMLReporter creates a new HTML reporter
func NewHTMLReporter() *HTMLReporter {
	tmpl := template.Must(template.New("report").Parse(htmlTemplate))
	return &HTMLReporter{
		template: tmpl,
	}
}

// GenerateReport generates an HTML report for the given test result
func (r *HTMLReporter) GenerateReport(result *testing.TestResult, outputPath string) error {
	// Create output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Prepare template data
	data := struct {
		Result       *testing.TestResult
		GeneratedAt  time.Time
		Duration     time.Duration
		PassCount    int
		FailCount    int
		WarningCount int
		ErrorCount   int
	}{
		Result:      result,
		GeneratedAt: time.Now(),
		Duration:    result.Duration,
	}

	// Count events by status
	for _, event := range result.Trace {
		switch event.Status {
		case "success":
			data.PassCount++
		case "failure", "fail":
			data.FailCount++
		case "warning":
			data.WarningCount++
		case "error":
			data.ErrorCount++
		}
	}

	// Execute template
	if err := r.template.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// GenerateSummaryReport generates a summary HTML report for multiple test results
func (r *HTMLReporter) GenerateSummaryReport(results []*testing.TestResult, outputPath string) error {
	// Create output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Calculate summary stats
	totalTests := len(results)
	passedTests := 0
	failedTests := 0
	totalDuration := time.Duration(0)
	totalEvents := 0

	for _, result := range results {
		if result.Passed {
			passedTests++
		} else {
			failedTests++
		}
		totalDuration += result.Duration
		totalEvents += len(result.Trace)
	}

	// Prepare template data
	data := struct {
		Results       []*testing.TestResult
		GeneratedAt   time.Time
		TotalTests    int
		PassedTests   int
		FailedTests   int
		TotalDuration time.Duration
		TotalEvents   int
		SuccessRate   float64
	}{
		Results:       results,
		GeneratedAt:   time.Now(),
		TotalTests:    totalTests,
		PassedTests:   passedTests,
		FailedTests:   failedTests,
		TotalDuration: totalDuration,
		TotalEvents:   totalEvents,
		SuccessRate:   float64(passedTests) / float64(totalTests) * 100,
	}

	// Execute template
	if err := r.template.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoMailServer Test Report - {{.Result.Name}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f5f5;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
        }

        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }

        .header p {
            font-size: 1.2em;
            opacity: 0.9;
        }

        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            padding: 30px;
            background: #f8f9fa;
            border-bottom: 1px solid #dee2e6;
        }

        .metric {
            text-align: center;
        }

        .metric h3 {
            font-size: 2em;
            margin-bottom: 5px;
        }

        .metric p {
            color: #6c757d;
            text-transform: uppercase;
            font-size: 0.9em;
            letter-spacing: 0.5px;
        }

        .status {
            padding: 20px 30px;
            text-align: center;
            font-size: 1.5em;
            font-weight: bold;
        }

        .status.pass {
            background-color: #d4edda;
            color: #155724;
        }

        .status.fail {
            background-color: #f8d7da;
            color: #721c24;
        }

        .content {
            padding: 30px;
        }

        .section {
            margin-bottom: 40px;
        }

        .section h2 {
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
            margin-bottom: 20px;
            color: #333;
        }

        .trace-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }

        .trace-table th,
        .trace-table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #dee2e6;
        }

        .trace-table th {
            background-color: #f8f9fa;
            font-weight: 600;
            color: #495057;
        }

        .trace-table tr:hover {
            background-color: #f8f9fa;
        }

        .status-badge {
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.8em;
            font-weight: bold;
            text-transform: uppercase;
        }

        .status-success {
            background-color: #d4edda;
            color: #155724;
        }

        .status-error {
            background-color: #f8d7da;
            color: #721c24;
        }

        .status-warning {
            background-color: #fff3cd;
            color: #856404;
        }

        .status-info {
            background-color: #d1ecf1;
            color: #0c5460;
        }

        .details {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 4px;
            margin-top: 10px;
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
        }

        .footer {
            text-align: center;
            padding: 20px;
            background: #f8f9fa;
            color: #6c757d;
            border-top: 1px solid #dee2e6;
        }

        @media (max-width: 768px) {
            .header {
                padding: 20px;
            }

            .header h1 {
                font-size: 2em;
            }

            .summary {
                grid-template-columns: repeat(2, 1fr);
                padding: 20px;
            }

            .content {
                padding: 20px;
            }

            .trace-table {
                font-size: 0.8em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Result.Name}}</h1>
            <p>{{.Result.Description}}</p>
        </div>

        <div class="summary">
            <div class="metric">
                <h3>{{if .Result.Passed}}✓{{else}}✗{{end}}</h3>
                <p>Status</p>
            </div>
            <div class="metric">
                <h3>{{.Duration}}</h3>
                <p>Duration</p>
            </div>
            <div class="metric">
                <h3>{{len .Result.Trace}}</h3>
                <p>Trace Events</p>
            </div>
            <div class="metric">
                <h3>{{len .Result.Errors}}</h3>
                <p>Errors</p>
            </div>
        </div>

        <div class="status {{if .Result.Passed}}pass{{else}}fail{{end}}">
            {{if .Result.Passed}}
                ✅ Test Passed Successfully
            {{else}}
                ❌ Test Failed - {{.Result.Summary}}
            {{end}}
        </div>

        <div class="content">
            {{if .Result.Errors}}
            <div class="section">
                <h2>Errors</h2>
                {{range .Result.Errors}}
                <div class="details">{{.}}</div>
                {{end}}
            </div>
            {{end}}

            <div class="section">
                <h2>Execution Trace</h2>
                <table class="trace-table">
                    <thead>
                        <tr>
                            <th>Time</th>
                            <th>Phase</th>
                            <th>Component</th>
                            <th>Action</th>
                            <th>Duration</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Result.Trace}}
                        <tr>
                            <td>{{.Timestamp.Format "15:04:05.000"}}</td>
                            <td>{{.Phase}}</td>
                            <td>{{.Component}}</td>
                            <td>{{.Action}}</td>
                            <td>{{.Duration}}</td>
                            <td>
                                <span class="status-badge status-{{.Status}}">
                                    {{.Status}}
                                </span>
                            </td>
                        </tr>
                        {{if .Details}}
                        <tr>
                            <td colspan="6">
                                <div class="details">
                                    {{range $key, $value := .Details}}
                                    <strong>{{$key}}:</strong> {{$value}}<br>
                                    {{end}}
                                </div>
                            </td>
                        </tr>
                        {{end}}
                        {{if .Error}}
                        <tr>
                            <td colspan="6">
                                <div class="details" style="background-color: #f8d7da; color: #721c24;">
                                    <strong>Error:</strong> {{.Error}}
                                </div>
                            </td>
                        </tr>
                        {{end}}
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>

        <div class="footer">
            <p>Report generated at {{.GeneratedAt.Format "2006-01-02 15:04:05"}} | GoMailServer Test Suite</p>
        </div>
    </div>
</body>
</html>`

// GenerateQuickReport creates a simple HTML report for immediate viewing
func (r *HTMLReporter) GenerateQuickReport(result *testing.TestResult) (string, error) {
	// Generate HTML content
	var htmlContent strings.Builder

	htmlContent.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Quick Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .pass { color: green; }
        .fail { color: red; }
        .summary { background: #f0f0f0; padding: 10px; border-radius: 5px; }
        .trace { margin: 10px 0; padding: 10px; border-left: 3px solid #ccc; }
        .success { border-left-color: green; }
        .error { border-left-color: red; }
    </style>
</head>
<body>
    <h1>GoMailServer Test Report</h1>
    <div class="summary">
        <h2>` + result.Name + `</h2>
        <p><strong>Status:</strong> `)

	if result.Passed {
		htmlContent.WriteString(`<span class="pass">PASSED</span>`)
	} else {
		htmlContent.WriteString(`<span class="fail">FAILED</span>`)
	}

	htmlContent.WriteString(`</p>
        <p><strong>Duration:</strong> ` + result.Duration.String() + `</p>
        <p><strong>Description:</strong> ` + result.Description + `</p>
    </div>

    <h3>Trace Events</h3>`)

	for _, event := range result.Trace {
		cssClass := "trace"
		if event.Status == "success" {
			cssClass += " success"
		} else if event.Status == "error" {
			cssClass += " error"
		}

		htmlContent.WriteString(`<div class="` + cssClass + `">
            <strong>` + event.Timestamp.Format("15:04:05.000") + `</strong> - ` +
			event.Phase + ` (` + event.Component + `) - ` + event.Action + `<br>
            <em>Duration:</em> ` + event.Duration.String() + ` | <em>Status:</em> ` + event.Status)

		if event.Error != "" {
			htmlContent.WriteString(` | <em>Error:</em> ` + event.Error)
		}

		htmlContent.WriteString(`</div>`)
	}

	htmlContent.WriteString(`</body></html>`)

	return htmlContent.String(), nil
}
