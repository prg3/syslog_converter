# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git gcc
WORKDIR $GOPATH/src/github.com/prg3/syslog_converter/app
ADD app .
RUN go get -d -v ; go build -o /syslog_converter

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /syslog_converter /app/
ENTRYPOINT ["./syslog_converter"]