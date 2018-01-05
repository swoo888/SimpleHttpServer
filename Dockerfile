FROM golang:1.9.2
WORKDIR /go/src/SpaceX
COPY . .
RUN go get ./... && go install .
CMD ["/go/bin/SpaceX"]