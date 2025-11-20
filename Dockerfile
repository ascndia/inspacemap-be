# --- Stage 1: Builder ---
FROM golang:1.21-alpine AS builder

# Install git dan dependency dasar (diperlukan untuk fetch go modules tertentu)
RUN apk add --no-cache git

WORKDIR /app

# Copy file dependensi dulu agar ter-cache layer-nya
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary
# CGO_ENABLED=0 memastikan binary statis (tidak butuh library C eksternal)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# --- Stage 2: Runner (Production Image) ---
FROM alpine:latest

# Install sertifikat SSL (penting jika aplikasi request ke HTTPS luar)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary dari Stage Builder
COPY --from=builder /app/main .

# (Opsional) Copy file .env jika ingin dibaca dari file, 
# tapi best practice di Docker adalah inject via environment variable
# COPY .env .

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]