FROM golang:1.17.6 as Builder
COPY . /go/src
WORKDIR /go/src
# set env for build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# build
RUN ["go", "install"]
RUN ["go", "build", "-o", "server"]

FROM alpine:latest
COPY --from=builder /go/src/server ./
ENV SERVER_ADDR=localhost:8081
ENV GODEBUG=http2debug=1
EXPOSE 8080
ENTRYPOINT [ "./server" ]
