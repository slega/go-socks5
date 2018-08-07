
#build stage
FROM golang:alpine AS stage1
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git g++ gcc make musl musl-dev 
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=stage1 /go/bin/app /app
ENTRYPOINT ./app
LABEL Name=go-socks5 Version=0.0.1
EXPOSE 8000 8080
