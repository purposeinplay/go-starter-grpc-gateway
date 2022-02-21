package psql_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	// nolint: revive // reports line to long
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql/psqldocker"

	// nolint: revive // reports line to long
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql/psqltest"

	// nolint: revive // reports line to long
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql/psqlutil"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/net"
)

const (
	usr      = "test"
	password = "pass"
	dbName   = "test"
)

var (
	port = "5432"
	dsn  = fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s h"+
			"ost=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		port,
	)
)

func TestMain(m *testing.M) {
	schema, err := psqlutil.ReadSchema()
	if err != nil {
		log.Println("err while reading schema:", err)
		os.Exit(1)
	}

	log.Println("Starting psql docker container...")

	portInt, err := net.GetFreePort()
	if err != nil {
		log.Println("err while getting free port:", err)
		os.Exit(1)
	}

	port := strconv.Itoa(portInt)

	dsn = fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s h"+
			"ost=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		port,
	)

	res, err := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithPort(port),
		psqldocker.WithContainerName("psql_tests"),
		psqldocker.WithSchema(schema),
	)
	if err != nil {
		log.Println("err while starting docker db:", err)
		os.Exit(1)
	}

	psqltest.Register(dsn)

	log.Println("Starting tests...")

	ret := m.Run()

	log.Println("Tests done. Teardown DB.")

	err = res.Close()
	if err != nil {
		log.Println("err while tearing down db:", err)
	}

	os.Exit(ret)
}
