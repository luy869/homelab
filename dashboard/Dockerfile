# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go backend with embedded frontend
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY backend/ .
COPY --from=frontend /frontend/dist ./static
RUN go build -mod=vendor -o dashboard .

# Stage 3: Minimal runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/dashboard .
EXPOSE 3000
CMD ["./dashboard"]
