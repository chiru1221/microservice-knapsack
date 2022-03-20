FROM golang:1.17.6 as builder
COPY . /go/src
WORKDIR /go/src
# set env for build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# build
RUN ["go", "install"]
WORKDIR /go/src/update
RUN ["go", "build", "-o", "server"]
WORKDIR /go/src/health
RUN ["go", "build", "-o", "healthz"]

FROM alpine:latest
COPY --from=builder /go/src/update/server /go/src/health/healthz ./
ENV SERVER_ADDR=localhost:8081
EXPOSE 8081
ENTRYPOINT [ "./server" ]
