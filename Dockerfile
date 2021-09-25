FROM golang
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && \
    go build -ldflags="-s -w" -o /pumpkin-pi

ENTRYPOINT ["/pumpkin-pi"]
