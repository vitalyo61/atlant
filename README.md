# atlant

run test:

```make test```

run integration test:

```make```

## make proto
```protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/atlant.proto```
