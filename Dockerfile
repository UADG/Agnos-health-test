# 1. Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# ก็อปปี้ไฟล์ go.mod และโหลด dependencies
COPY go.mod go.sum ./
RUN go mod download

# ก็อปปี้โค้ดทั้งหมดและ Build เป็นไฟล์ชื่อ main
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go
# 2. Run Stage (ใช้ Alpine เปล่าๆ เพื่อให้ Image มีขนาดเล็กที่สุด)
FROM alpine:latest
WORKDIR /app

# เพิ่มบรรทัดนี้เพื่อติดตั้ง Timezone Data ให้ Alpine รู้จัก Asia/Bangkok
RUN apk add --no-cache tzdata

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]