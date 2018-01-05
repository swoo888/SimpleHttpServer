FROM golang:1.9.2
WORKDIR /go/src/SimpleHttpServer
COPY . .
RUN go get ./... && go install .
CMD ["/bin/bash", "-c", "/go/bin/SimpleHttpServer -db postgres -host postgresdb -pass \
    ${POSTGRES_PASSWORD} -user postgres"]