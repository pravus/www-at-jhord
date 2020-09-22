FROM alpine:3 as builder

RUN apk --no-cache update \
 && apk --no-cache upgrade \
 && apk --no-cache add git go \
 && mkdir -p /usr/src/www-at-jhord

COPY *.go /usr/src/www-at-jhord/

RUN cd /usr/src/www-at-jhord \
 && go get -d \
 && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"'


FROM scratch

COPY --from=builder /usr/src/www-at-jhord/www-at-jhord /
COPY resume /resume/

ENTRYPOINT ["/www-at-jhord"]
