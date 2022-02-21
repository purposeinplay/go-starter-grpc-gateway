package command

import (
	"context"
	"fmt"
)

// ReportError contains the data needed
// to report an error to an external service.
type ReportError struct {
	Err error
}

// ReportErrorHandler contains the dependencies for reporting
// an error.
type ReportErrorHandler struct {
	reportService ReportService
}

// NewReportErrorHandler returns a new ReportErrorHandler.
func NewReportErrorHandler(
	reportService ReportService,
) ReportErrorHandler {
	return ReportErrorHandler{
		reportService: reportService,
	}
}

// Handle reports an error to an external service.
func (h ReportErrorHandler) Handle(ctx context.Context, cmd ReportError) error {
	err := h.reportService.ReportError(ctx, cmd.Err)
	if err != nil {
		return fmt.Errorf("report: %w", err)
	}

	return nil
}
