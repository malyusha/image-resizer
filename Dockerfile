FROM golang:1.10-alpine

RUN apk add --no-cache git build-base curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/malyusha/image-resizer

ADD ./Gopkg.lock ./Gopkg.lock
ADD ./Gopkg.toml ./Gopkg.toml
RUN dep ensure -vendor-only

ADD . .
ADD ./start.sh ./start.sh

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./bin/resizer ./cmd/resize

ENV PATH="/go/src/github.com/malyusha/image-resizer/bin:${PATH}"

EXPOSE 9898

ENTRYPOINT ["./start.sh"]