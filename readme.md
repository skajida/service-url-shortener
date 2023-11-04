# URL Shortener service

## Dependencies

### Deployment

[Podman](https://podman.io/docs/installation) utility was used to deploy the service. It has better integration with Kubernetes than the mainstream container tool – Docker, but still uses compliant container formats.

### Code generation

[protobuf](https://protobuf.dev/) used for [gRPC wrappers](/internal/pb/) and [swagger file](/api/url-shortener.swagger.json) code generation.

To install it on rhel-like linux distributions use this command:
```bash
# enterprise linux 9
sudo dnf install protobuf-devel
```

The following version of the utility was used:
```console
$ protoc --version
libprotoc 25.0
```

`protoc` dependencies:
* [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go)
* [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)
* [protoc-gen-grpc-gateway](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway)
* [protoc-gen-openapiv2](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2)

gRPC wrappers and swagger file were generated this way:
```bash
protoc -I proto --go_out=./internal/pb --go_opt paths=source_relative --go-grpc_out=./internal/pb --go-grpc_opt paths=source_relative --grpc-gateway_out=./internal/pb --grpc-gateway_opt paths=source_relative --openapiv2_out=./api proto/url-shortener.proto
```

Mocks was generated from `contract.go` files with [mockgen](https://pkg.go.dev/go.uber.org/mock/mockgen) utility to write unit tests, the result is two files:
```console
$ find . -name contract.go
./internal/service/contract.go
./internal/controller/contract.go
$ # they both have go-generate-related command at the top of the file
$ head -n1 ./internal/service/contract.go
//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
$ go generate ./...
$ find . -name mocks_test.go
./internal/service/mocks_test.go
./internal/controller/mocks_test.go
```

## Build

Example of service building command:
```bash
go build -o bin/service-url-shortener ./cmd
```

But to simplify future deployment better build container images:
```bash
podman build . -f build/db/Containerfile -t service-url-shortener-db
podman build . -f build/app/Containerfile -t service-url-shortener-app
```

## Deploy

The service looks up for an required non-empty environment variable `STORAGE_TYPE` to select the corresponding one storage type.

Practically it is configured by selecting the corresponding k8s pod configuration. Both of them are in [`/deploy`](/deploy/) directory.

To deploy service with postgres-like storage run the following command:
```bash
podman kube play --configmap deploy/app-configmap.yml deploy/postgres-pod.yml
```

For the in-memory-like storage:
```bash
podman kube play --configmap deploy/app-configmap.yml deploy/in-memory-pod.yml
```

The following command will stop and remove all related containers. Use `--force` flag for postgres-like deployment to prune related to database persistent volume.
```bash
podman kube down deploy/STORAGE-TYPE-pod.yml
```

## Connect

There is generated [swagger file](/api/url-shortener.swagger.json) describing how to send requests to service. Some examples:
* `POST`
    ```bash
    curl --location 'localhost:8081/reduce' --header 'Accept: application/json' --header 'Content-Type: text/plain' --data '{"originUrl": "https://www.ozon.ru/"}'
    ```
* `GET`:
    ```bash
    curl --location 'localhost:8081/get/tUaVM2YjHq' --header 'Accept: application/json'
    ```

## Test

There are 3 modules that have been covered by unit-tests, to run them run:
```bash
go test ./internal/controller
go test ./internal/entity
go test ./internal/service
```
