package mocks

import (
	"context"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/stretchr/testify/mock"
)

var _ command.ReportService = (*ReportService)(nil)

type ReportService struct {
	mock.Mock
}

func (m *ReportService) ReportError(ctx context.Context, err error) error {
	args := m.Called(ctx, err)

	return args.Error(0)
}
