package psqldocker

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// NewContainer starts a new psql database in a docker container.
func NewContainer(
	user, password, dbName string,
	opts ...Option,
) (*dockertest.Resource, error) {
	var options options

	for _, o := range opts {
		o.apply(&options)
	}

	port, portBindings := ports(options)

	pool, err := pool(options)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	// create run options
	dockerRunOptions := &dockertest.RunOptions{
		Name:         options.containerName,
		Repository:   "postgres",
		Tag:          "latest",
		PortBindings: portBindings,
		Env:          envVars(user, password, dbName),
	}

	res, err := startContainer(pool, dockerRunOptions, func() error {
		return pingDB(user, password, dbName, port)
	})
	if err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	err = schema(user, password, dbName, port, options)
	if err != nil {
		_ = res.Close()

		return nil, fmt.Errorf("execute schema: %w", err)
	}

	return res, nil
}

func startContainer(
	pool *dockertest.Pool,
	runOptions *dockertest.RunOptions,
	retryFunc func() error,
) (*dockertest.Resource, error) {
	res, err := pool.RunWithOptions(
		runOptions,
		func(config *docker.HostConfig) {
			config.AutoRemove = true
		},
	)
	if err != nil {
		return nil, fmt.Errorf("docker run: %w", err)
	}

	err = pool.Retry(retryFunc)
	if err != nil {
		_ = res.Close()

		return nil, fmt.Errorf("ping node: %w", err)
	}

	return res, nil
}

func pool(opts options) (*dockertest.Pool, error) {
	pool := opts.pool

	if pool == nil {
		p, err := dockertest.NewPool("")
		if err != nil {
			return nil, fmt.Errorf("new pool: %w", err)
		}

		p.MaxWait = 20 * time.Second

		pool = p
	}

	return pool, nil
}

func ports(opts options) (
	port string,
	portBindings map[docker.Port][]docker.PortBinding,
) {
	port = "5432"

	if opts.port != "" {
		port = opts.port
	}

	pB := map[docker.Port][]docker.PortBinding{
		docker.Port("5432/tcp"): {{HostIP: "0.0.0.0", HostPort: port}},
	}

	return port, pB
}

func envVars(user, password, dbName string) []string {
	const envVars = 3

	env := make([]string, 0, envVars)

	env = append(env,
		fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		fmt.Sprintf("POSTGRES_USER=%s", user),
		fmt.Sprintf("POSTGRES_DB=%s", dbName),
	)

	return env
}

func pingDB(user, password, dbName, port string) error {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s h"+
			"ost=localhost "+
			"port=%s "+
			"sslmode=disable",
		user,
		password,
		dbName,
		port))
	if err != nil {
		return err
	}

	defer func() {
		_ = db.Close()
	}()

	return db.Ping()
}

func schema(
	user,
	password,
	dbName,
	port string,
	options options,
) error {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s h"+
			"ost=localhost "+
			"port=%s "+
			"sslmode=disable",
		user,
		password,
		dbName,
		port))
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	defer func() {
		_ = db.Close()
	}()

	if options.schema != "" {
		_, err = db.Exec(options.schema)
		if err != nil {
			return fmt.Errorf("execute schema: %w", err)
		}
	}

	return nil
}
