# Description

This project builds the `@jhord` sub-site of [www.carbon.cc](www.carbon.cc).

# Building Locally

## Bare Metal

```bash
go build && ./www-at-jhord
```

## With Docker

```bash
docker build -t www-at-jhord:latest . && docker run --name=at-jhord-www --rm -it --network=host www-at-jhord:latest
```
