# Description

This project builds the `@jhord` sub-site of [www.carbon.cc](www.carbon.cc).

# Development

## Bare Metal

In order to develop locally on bare metal you will first need to create the registry stubs from the .proto file.

```bash
mkdir -p registry && protoc --gofast_out=plugins=grpc:registry registry.proto
```

## Using Docker

In order to build and test locally you will need to create a local Docker network for both the HTTP and GRPC containers.
The following command creates the `www` network:


```bash
docker network create www
```

To start the HTTP container:

```bash
docker build -f Dockerfile.http -t www-at-jhord-http:latest . && docker run --name=at-jhord-www-http --rm -it --publish=8000:8000 --network=www www-at-jhord-http:latest
```

To start the GRPC container:

```bash
docker build -f Dockerfile.grpc -t www-at-jhord-grpc:latest . && docker run --name=www-at-jhord-grpc --rm -it --network=www www-at-jhord-grpc:latest
```

## Links

* [Home](http://localhost:8000/)
* [Resume](http://localhost:8000/resume)
* [Visits](http://localhost:8000/visits)
