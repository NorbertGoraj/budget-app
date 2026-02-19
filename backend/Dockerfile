FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /budget-api . && \
    CGO_ENABLED=0 go build -o /migrate ./cmd/migrate

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY --from=builder /budget-api /budget-api
COPY --from=builder /migrate /migrate
EXPOSE 8080
CMD ["/budget-api"]
