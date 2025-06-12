FROM golang:1.24.3-alpine AS builder
WORKDIR /src
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/*.go src/root.html ./
RUN go build -o training-application

FROM alpine:3.21.3
WORKDIR /app
COPY --from=builder /src/training-application /app/training-application
COPY training-application.conf .
ENTRYPOINT [ "./training-application" ]
