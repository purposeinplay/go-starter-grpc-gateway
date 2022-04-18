package grpc

import (
	"github.com/purposeinplay/go-commons/psqltest"
	"log"
	"os"
	"testing"

	// revive:disable-next-line:line-length-limit
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.LoadTestConfig("../../../config.test.yaml")

	if err != nil {
		log.Fatalf("error loading config: %+v", err)
	}

	psqltest.Register(dsn)

	exitVal := m.Run()

	os.Exit(exitVal)
}
