## Potz Holzöpfel und Zipfelchape

Tri Tra Trulla La!

## Features

* Hosts an embedded copy of Caasperli on port 8080
* Available as minimal container image (`FROM scratch`)
* Prometheus endpoint availabe (`/metrics`)
* Supports tracing
* May use MongoDB as persistance layer to display a simple web counter (optional)

## Usage

Run an `potz-holzoepfel-und-zipfelchape` binary from the [releases page](https://github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/releases) or use
the container available at `docker.pkg.github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/app`.

```bash
docker run --rm -ti -p 8080:8080 docker.pkg.github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/app
```

Get page contents:
```bash
curl localhost:8080
```

Get metrics:
```bash
curl localhost:8080/metrics
```

### Configure Jaeger tracing middleware

Caasperli is instrumented using [`jaeger-client-go`](https://github.com/jaegertracing/jaeger-client-go). 

The instrumentation layer is configured via environment variables. The following example shows how to
send spans to the jaeger collector over TCP while sampling each request to ensure you get some test
data in Jaeger ASAP.

```bash
export JAEGER_ENDPOINT=http://jaeger.example.com:14268/api/traces
export JAEGER_SAMPLER_TYPE=const
export JAEGER_SAMPLER_PARAM=1

potz-holzoepfel-und-zipfelchape
```

If you enable persistence as described below then MongoDB operations will be added
to the default requests spans.

### Running with enabled persistence

You can run Caasperli with an optional persistance layer and it will display a
simple hit counter with data from MongoDB.

```bash
potz-holzoepfel-und-zipfelchape \
  --persistence
  --mongodb-uri="mongodb://root:hunter2@localhost:27017"
```

### Persistence on Cloud Foundry

If you are running on Cloud Foundry you can activate the persistence layer by
attaching a MongoDB service called `mongodb` to the application.

```
cf push caasperli --no-start
cf create-service mongodb <plan> mongodb
cf bind-service caasperli mongodb
cf start caasperli
```

Caasperli will detect that it is running on Cloud Foundry and automatically
activate it's persistance layer if a service called `mongodb` is bound to the
app.

### Deploying with Waypoint

```
waypoint up
````

## Development

Run a server.

```bash
go run main.go
```

Build a local copy of the server.

```bash
go build
```

### MongoDB persistance layer

There is a `docker-compose.yml` to run a local MongoDB instance that works
with the default settings.

```bash
docker-compose up -d
go run main.go --persistence
```

### Release Process

Create a git tag and push it to this repo or use the git web ui.

## License

This application is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

It is heavily based on [enricofoltran/simple-go-server](https://github.com/enricofoltran/simple-go-server) which is licensed under the MIT license.
