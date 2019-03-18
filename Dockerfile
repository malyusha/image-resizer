FROM golang:1.12-alpine as builder
ARG SOURCE="/go/src/github.com/malyusha/image-resizer"

ADD . ${SOURCE}
WORKDIR ${SOURCE}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./go/bin/resizer ./cmd/resizer

FROM debian
COPY --from=builder /go/bin/resizer /bin/image-resizer

EXPOSE 9898

ENTRYPOINT ["image-resizer"]