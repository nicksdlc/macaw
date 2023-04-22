# macaw
Macaw is a tool to mock remote service. It can work in "sender" and "responder" modes.

In respective mode it will either send a defined amount of requests or respond to request with specified response.

## Configuration

Macaw is configured via config.yml. Example is repository root.

## Run

Macaw can be executed either by
```
go run main.go
```

Or by building it to executable.

Features:
- Macaw metalang
- Fluent API for macaw to configure 
- Full RabbitMQ protocol mock
