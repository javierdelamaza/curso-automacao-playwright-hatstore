FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY backend backend
COPY frontend frontend
WORKDIR /app/backend
RUN go build -o /app/app

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY frontend /app/frontend
EXPOSE 8080
CMD ["/app/app"]