FROM golang:1.9.2
WORKDIR /go/src/SimpleHttpServer
COPY . .
RUN go get ./... && go install .
CMD ["/go/bin/SimpleHttpServer"]