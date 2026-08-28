package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	deliveryHttp "agnos-assignment/delivery/http"
	"agnos-assignment/delivery/http/middleware"
	repoPostgres "agnos-assignment/repository/postgres"
	repoExternal "agnos-assignment/repository/external"
	"agnos-assignment/usecase"
)

func main() {
	// 1. กำหนดค่า Configuration (ในของจริงควรดึงจาก .env หรือ Environment Variables)
	dsn := "host=postgres_db user=postgres password=postgres dbname=agnos_db port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	jwtSecret := "your-super-secret-key-for-jwt" // ต้องตรงกับที่ใช้ใน Docker/Server

	// 2. เชื่อมต่อ Database ด้วย GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// 3. เริ่มต้น Web Framework (Gin)
	router := gin.Default()

	// 4. สร้าง Middleware
	jwtAuthMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	// ==========================================
	// 5. Dependency Injection (DI) - ประกอบร่าง
	// ==========================================

	// Layer 3: Repositories & External API Clients
	// สร้างตัวจัดการฐานข้อมูลและ API ภายนอก
	hospitalRepo := repoPostgres.NewHospitalRepository(db)
	staffRepo := repoPostgres.NewStaffRepository(db)
	patientRepo := repoPostgres.NewPatientRepository(db)
	
	// สมมติว่าคุณสร้าง HTTP Client สำหรับเรียก Hospital A API ไว้ในโฟลเดอร์ repository/external
	hisClient := repoExternal.NewHISClient("https://hospital-a.api.co.th")

	// Layer 2: Usecases
	// นำ Repositories ฉีดเข้าไปใน Usecase
	staffUsecase := usecase.NewStaffUsecase(staffRepo, hospitalRepo, jwtSecret)
	patientUsecase := usecase.NewPatientUsecase(patientRepo, hisClient)

	// Layer 4: Handlers (Delivery)
	// นำ Usecase และ Middleware ฉีดเข้าไปใน Gin Handler พร้อมกับผูก Route
	deliveryHttp.NewStaffHandler(router, staffUsecase)
	deliveryHttp.NewPatientHandler(router, patientUsecase, jwtAuthMiddleware)

	// ==========================================
	// 6. รันเซิร์ฟเวอร์
	// ==========================================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}