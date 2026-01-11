package checks

import (
	"context"

	"github.com/btafoya/gomailserver/internal/testing/types"
)

type Check interface {
	Name() string
	Description() string
	Category() types.Category
	Severity() types.Severity
	Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error)
}
