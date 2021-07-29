GoStarter
==========

> This starter kit is designed to get you up and running with a project structure optimized for developing RESTful API services in Go. It is an opinionated Go starter kit built on top of Chi, using battle tested libraries proven to provide a good foundation for a project written in Golang.

### Prerequisites

The codebase requires these development tools:

* Go compiler and runtime: 1.15.2 or greater.
* Protobuf: 3.14.0 or greater.
* Docker Engine: 19.0.0 or greater.

You must install the Go plugins for the Protobuf toolchain:

```
env GO111MODULE=off go install "github.com/golang/protobuf/protoc-gen-go"
```

For the latest instructions see the [official documentation](https://developers.google.com/protocol-buffers/docs/gotutorial).

### Go Dependencies

The project uses Go modules which should be vendored:

```shell
env GO111MODULE=on GOPRIVATE="github.com" go mod vendor
```

You can regenerate the Protobuf Go generated stubs:

```
make proto
```

### Configuration

You may configure GoStarter using either a configuration file named .env, environment variables, or a combination of both. Environment variables are prefixed with GOSTARTER, and will always have precedence over values provided via file.

#### Server
```properties
SERVER_ADDRESS: localhost
SERVER_PORT: 7350
```

`ADDRESS` - `string`

Hostname to listen on.

`PORT` - `number`

Port number to listen on. Defaults to `7350`.

#### Database

```properties
DB_DRIVER: postgres
DB_HOST: dbhost
DB_USER: dbuser
DB_PASSWORD: dbpassword
DB_NAME: dbname
```

**Migrations Note** Migrations are not applied automatically, so you will need to run them after you've built GoStarter.
* If built locally: `./gostarter migrate`
* Using Docker: `docker run --rm gostarter gostarter migrate`

### Start in Development

The recommended workflow is to use Docker and the compose file to build and run the service and resources.

```shell
docker-compose -f docker-compose.dev.yml up
```

__Hot Reloading:__ GoStarter uses Reflex in development for hot reloading.

### Start In Production

```shell
make image
```

now run the image with
```shell
docker run -it go-user
```

### Run Tests
Running the tests locally requires a valid database connection. Configure `config.test.yaml` with the appropriate values.

```shell
make test
```