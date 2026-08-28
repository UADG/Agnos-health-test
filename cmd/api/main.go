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
	dsn := "host=postgres_db user=postgres password=postgres dbname=agnos_db port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	jwtSecret := "your-super-secret-key-for-jwt" 

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	router := gin.Default()

	jwtAuthMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	hospitalRepo := repoPostgres.NewHospitalRepository(db)
	staffRepo := repoPostgres.NewStaffRepository(db)
	patientRepo := repoPostgres.NewPatientRepository(db)
	
	hisClient := repoExternal.NewHISClient("https://hospital-a.api.co.th")

	staffUsecase := usecase.NewStaffUsecase(staffRepo, hospitalRepo, jwtSecret)
	patientUsecase := usecase.NewPatientUsecase(patientRepo, hisClient)

	deliveryHttp.NewStaffHandler(router, staffUsecase)
	deliveryHttp.NewPatientHandler(router, patientUsecase, jwtAuthMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}