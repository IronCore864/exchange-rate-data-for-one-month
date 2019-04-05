FROM golang:alpine AS build-env
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache git
WORKDIR $GOPATH/src/github.com/ironcore864/exchange-rate-data-for-one-month
COPY . .
RUN go get -u github.com/go-redis/redis && go build -o goapp

FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=build-env /go/src/github.com/ironcore864/exchange-rate-data-for-one-month/goapp /app/
ENTRYPOINT ./goapp
