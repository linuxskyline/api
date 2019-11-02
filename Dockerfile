FROM golang:alpine as build

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/api

FROM alpine:latest

COPY --from=build /go/bin/api /api

EXPOSE 80
ENTRYPOINT ["/api"]
