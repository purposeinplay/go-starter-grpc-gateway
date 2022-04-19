package grpc_test

import (
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
	"testing"
)

var testID = uuid.New()

func TestMain(m *testing.M) {
	revert := uuid.SetTestUUID(testID)

	defer revert()

	m.Run()
}
