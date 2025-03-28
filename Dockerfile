FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN go build -o training-application

FROM alpine:3.21.3
WORKDIR /app
COPY --from=builder /src/training-application /app/training-application
COPY conf/app.conf ./conf/
ENTRYPOINT [ "./training-application" ]
