package scenarios

import (
	"context"

	"github.com/btafoya/gomailserver/internal/testing"
)

// TestScenario defines the interface for test scenarios
type TestScenario interface {
	// Name returns the scenario name
	Name() string

	// Description returns a description of the scenario
	Description() string

	// Setup prepares the scenario for execution
	Setup(ctx context.Context) error

	// Execute runs the main scenario logic
	Execute(ctx context.Context) error

	// Verify checks that the scenario executed correctly
	Verify(ctx context.Context) error

	// Cleanup cleans up after the scenario
	Cleanup(ctx context.Context) error
}

// BaseScenario provides common functionality for test scenarios
type BaseScenario struct {
	name        string
	description string
	tracer      *testing.TraceCollector
}

// NewBaseScenario creates a new base scenario
func NewBaseScenario(name, description string, tracer *testing.TraceCollector) *BaseScenario {
	return &BaseScenario{
		name:        name,
		description: description,
		tracer:      tracer,
	}
}

// Name returns the scenario name
func (b *BaseScenario) Name() string {
	return b.name
}

// Description returns the scenario description
func (b *BaseScenario) Description() string {
	return b.description
}

// Setup provides default setup implementation (does nothing)
func (b *BaseScenario) Setup(ctx context.Context) error {
	return nil
}

// Cleanup provides default cleanup implementation (does nothing)
func (b *BaseScenario) Cleanup(ctx context.Context) error {
	return nil
}

// Trace returns the trace collector
func (b *BaseScenario) Trace() *testing.TraceCollector {
	return b.tracer
}
