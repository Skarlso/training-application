FROM golang:1.24.3-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o training-application

FROM alpine:3.21.3
WORKDIR /app
COPY --from=builder /src/training-application /app/training-application
COPY conf/app.conf ./conf/
ENTRYPOINT [ "./training-application" ]
