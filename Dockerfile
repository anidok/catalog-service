# Use official Golang image as build stage
FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o catalog-service ./cmd/api/main.go

# Use minimal base image for runtime
FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /app/catalog-service /app/catalog-service
COPY application.yaml /app/application.yaml

EXPOSE 4000

ENTRYPOINT ["/app/catalog-service"]
