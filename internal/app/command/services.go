package command

import (
	"context"
)

// ReportService is used to report outstanding behaviour.
type ReportService interface {
	ReportError(ctx context.Context, err error) error
}
