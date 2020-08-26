## Potz Holzöpfel und Zipfelchape

Tri Tra Trulla La!

## Features

* Hosts an embedded copy of Caasperli on port 8080
* Available as minimal container image (`FROM scratch`)
* Prometheus endpoint availabe (`/metrics`)
* Supports tracing

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

## Development

### Statik
Regenerate `statik/` dir after changing `public/` dir.

```bash
go generate
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
