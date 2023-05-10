# macaw
Macaw is a tool to mock remote service. It can work in "sender" and "responder" modes.

In respective mode it will either send a defined number of requests or respond to request with specified response.

## Configuration

Macaw is configured via config.yml. Example is repository root.

## Run

Macaw can be executed either by
```
go run main.go
```

Or by building it to executable.

## Configuration

### :godmode: Modes

Currently two modes are supported:
- Sender
- Responder

They are defined in the section of the config named 
```
Mode: receiver
```

### :electric_plug: Connectors 

To send and receive messages macaw should be configuraed with proper connectors.
At the moment HTTP and RabbitMQ. The respective configurations:

```
HTTP:
  Serve: 
    Host: /test
    Port: 1234
```

RabbitMQ:
```
Rabbit:
  Host: localhost
  Port: 5672
  User: guest
  Password: guest
  Exchanges: 
    - Exchange: 
      Name: ex.response
      RoutingKey: "responses.#"
    - Exchange: 
      Name: test-sender-requests
  Queues: 
    - Queue:
      Name: q.response
      Args:
        - "x-message-ttl": 300000
        - "x-queue-mode": "lazy"
    - Queue:
      Name: q.test
  ConnectionRetry:
    ElapsedTime: 2
    Interval: 20
```

### :envelope: Messages

#### Reponse

Configuration of reponse includes:
- **Request** - to which request response should be provided:
  - **To** - where to listen to request
  - **Matchers** - to which request to reply. More details on Matchers below
- **Body** - either file or string in the config. Defines template to use for the response
- **Options** - details like quantity and so on.

#### Request

Configuration of request:
- **To** - where to send the request
- **Body** - same as in response - body for the request
- **Options** - same as in response

For more details - refer to the _config_concept.yml_ in the repository root.

## Templates

Macaw uses go templates with some additions to make responses and requests more configurable.
The example of the response with templates is:

```
{
    "version": "1.0",
    "eventTimeUtc": "{{.Date}}",
    "status": "Success",
    "reason": "Retrieved",
    "requestId": {{.FromRequest "requestId"}},
    "policyEventTimeUtc": "{{.FromRequest "eventTimeUtc"}}",
    "bulkId": {{.Number "incremental"}},
    "bulksCount": {{.Quantity}} 
}
```

To define the placeholder use {{}} with the dot inside before the required placeholder.
At the moment following placeholders are supported:
- Number
- Date
- String
- FromRequest
- Quantity

Some of the placeholders support additional parameters, like _incremental_ for numbers - which will increment the number in the sequence of responses.