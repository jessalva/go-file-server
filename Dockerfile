FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app/
COPY . .
RUN go get ./...
RUN go build -o ./bin/GoFileServer ./cmd/GoFile/main.go

FROM alpine:3.9
WORKDIR /go/bin
COPY --from=build /go/src/app/bin/GoFileServer /go/bin/GoFileServer
EXPOSE 8090
ENTRYPOINT /go/bin/GoFileServer