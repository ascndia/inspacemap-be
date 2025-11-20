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

# 1. Build Binary SERVER
# CGO_ENABLED=0 memastikan binary statis (tidak butuh library C eksternal)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 2. Build Binary SEEDER
# Kita build juga seeder-nya agar bisa dijalankan di container production
RUN CGO_ENABLED=0 GOOS=linux go build -o seeder ./cmd/seeder/main.go

# --- Stage 2: Runner (Production Image) ---
FROM alpine:latest

# Install sertifikat SSL (penting jika aplikasi request ke HTTPS luar/S3 AWS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy Binary Server dari Builder
COPY --from=builder /app/main .

# Copy Binary Seeder dari Builder
COPY --from=builder /app/seeder .

# Copy Script Entrypoint
COPY entrypoint.sh .

# Beri izin eksekusi ke script (PENTING!)
RUN chmod +x entrypoint.sh

# Expose port aplikasi
EXPOSE 8080

# Ganti CMD untuk menjalankan script wrapper, bukan langsung main
# Script ini akan menjalankan seeder dulu, baru server
CMD ["./entrypoint.sh"]