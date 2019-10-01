FROM golang:1.13 AS build
COPY . /go/src/github.com/pirsquareff/dksweeper
WORKDIR /go/src/github.com/pirsquareff/dksweeper
RUN go get && CGO_ENABLED=0 go build

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/pirsquareff/dksweeper/dksweeper /dksweeper
ENTRYPOINT ["/dksweeper"]
