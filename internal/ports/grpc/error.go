package grpc

// and to use the google.golang.org one,
// but I don't know how to change this repo in the generate go files.
import (
	"context"
	"fmt"
	"github.com/purposeinplay/go-starter-grpc-gateway/apigrpc/v1"

	"github.com/golang/protobuf/proto"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) handleErr(err error) error {
	var (
		applicationError *errors.Error
		details          proto.Message
		grpcStatus       *status.Status
	)

	// In order to preserve space it would be better
	// to only log internal errors.
	s.logErr(err)

	// Check if the error is an application error or an
	// internal error
	switch {
	// If the error is an application error prepare the grpc
	// response.
	case errors.As(err, &applicationError):
		// Convert the application error type to a GRPC status.
		grpcStatus = errorToGRPCStatus(applicationError)

		// Convert the application error details to a grpc ErrorResponse
		// message.
		details = errorDetailsToGRPCDetails(applicationError)

	// If the error is an internal error, report it to an external
	// service.
	default:
		// Report the error to an external service
		reportErr := s.reportErr(err)
		if err != nil {
			// Log the error received from the external service
			// and continue execution.
			s.logErr(fmt.Errorf("report: %w", reportErr))
		}

		// Create a GRPC status that doesn't leak any information
		// about the internal error.
		grpcStatus = status.New(codes.Internal, "internal error.")
	}

	// Check if the details message is not nil
	// and attach it to the grpc status.
	if details != nil {
		grpcStatus, err = grpcStatus.WithDetails(details)
		if err != nil {
			// if attaching details to the grpc status
			// log and report the error.
			s.logErr(err)

			reportErr := s.reportErr(err)
			if reportErr != nil {
				s.logErr(fmt.Errorf("report: %w", err))
			}
		}
	}

	// Return the grpc Status as an immutable error.
	return grpcStatus.Err()
}

// errorToGRPCStatus converts an application defined error type
// to a grpc canonical error code.
func errorToGRPCStatus(err *errors.Error) *status.Status {
	var code codes.Code

	switch err.Type() {
	case errors.ErrorTypeInvalid:
		code = codes.InvalidArgument

	case errors.ErrorTypeNotFound:
		code = codes.NotFound

	default:
		code = codes.Unknown
	}

	return status.New(code, err.Message())
}

// errorDetailsToGRPCDetails checks if error details are attached to an
// application errors and converts them to a grpc message.
func errorDetailsToGRPCDetails(d *errors.Error) *startergrpc.ErrorResponse {
	details, ok := d.Details()
	if !ok {
		return nil
	}

	var code startergrpc.ErrorResponse_ErrorCode

	switch details.ApplicationErrorCode() {
	case errors.InternalErrorCodeNotEnoughBalance:
		code = startergrpc.ErrorResponse_ERROR_CODE_NOT_ENOUGH_BALANCE
	case errors.InternalErrorCodeEmailNotProvided:
		code = startergrpc.ErrorResponse_ERROR_CODE_EMAIL_NOT_PROVIDED
	}

	return &startergrpc.ErrorResponse{
		ErrorCode: code,
		Message:   details.Message(),
	}
}

// logErr logs the error.
func (s *Server) logErr(err error) {
	s.logger.Error("grpc", zap.Error(err))
}

// reportErr reports an error to an external service.
func (s *Server) reportErr(err error) error {
	reportErr := s.app.Commands.ReportError.Handle(
		context.Background(),
		command.ReportError{Err: err},
	)
	if reportErr != nil {
		return fmt.Errorf("report error command: %w", err)
	}

	return nil
}

func (s *Server) handlePanicRecover(p interface{}) error {
	s.logPanic(p)

	return status.Error(codes.Internal, "internal error.")
}

func (s *Server) logPanic(p interface{}) {
	s.logger.Error("grpc panic", zap.Any("cause", p))
}
