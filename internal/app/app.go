package app

import (
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
)

// Application represents the actions that can be
// called on the application.
type Application struct {
	Commands Commands
	Queries  Queries
}

// Commands represents the commands available in the application.
type Commands struct {
	CreateUser  command.CreateUserHandler
	ReportError command.ReportErrorHandler
}

// Queries represents the queries available in the application.
type Queries struct {
	FindUsers query.FindUsersHandler
	UserByID  query.UserByIDHandler
}
