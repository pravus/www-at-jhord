# Description

This project builds the `@jhord` sub-site of [www.carbon.cc](www.carbon.cc).

# Building Locally

This project can be built and run locally either as a stand-alone application or inside a Docker container.
By default the application accepts HTTP traffic on port 8000.
This can be changed by setting the `HTTP_BIND` environment variable prior to invocation.

## Stand-alone

```bash
go build && HTTP_BIND=:8000 ./www-at-jhord
```

## With Docker

```bash
docker build -t www-at-jhord:latest . && docker run --name=at-jhord-www --rm -it --publish=8000:8000 www-at-jhord:latest
```

## Links

* [Resume](http://localhost:8000/resume)
