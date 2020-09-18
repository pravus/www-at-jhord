FROM alpine:3 as builder

RUN apk --no-cache update \
 && apk --no-cache upgrade \
 && apk --no-cache add git go \
 && mkdir -p /usr/src/at-jhord

COPY *.go /usr/src/at-jhord/

RUN cd /usr/src/at-jhord \
 && go get -d \
 && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"'


FROM scratch

COPY --from=builder /usr/src/at-jhord/at-jhord /
COPY resume /resume/

ENTRYPOINT ["/at-jhord"]
