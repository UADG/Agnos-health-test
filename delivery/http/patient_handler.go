package http

import (
	"net/http"

	"agnos-assignment/domain"
	"github.com/gin-gonic/gin"
)

type patientHandler struct {
	patientUsecase domain.PatientUsecase
}

func NewPatientHandler(r *gin.Engine, pu domain.PatientUsecase, jwtMiddleware gin.HandlerFunc) {
	handler := &patientHandler{patientUsecase: pu}

	// สร้าง Group API และใส่ Middleware เพื่อบังคับ Login
	protected := r.Group("/patient")
	protected.Use(jwtMiddleware)
	{
		protected.GET("/search", handler.Search)
	}
}

func (h *patientHandler) Search(c *gin.Context) {
	// 1. ดึง hospital_id ที่ Middleware ฝากเอาไว้ (ป้องกันพนักงานค้นหาข้ามโรงพยาบาล)
	hospitalID := c.MustGet("hospital_id").(string)

	// 2. รับค่าจาก Query Parameters (?national_id=xxx&first_name=yyy)
	criteria := domain.PatientSearchCriteria{
		NationalID:  c.Query("national_id"),
		PassportID:  c.Query("passport_id"),
		FirstName:   c.Query("first_name"),
		LastName:    c.Query("last_name"),
		DateOfBirth: c.Query("date_of_birth"),
		PhoneNumber: c.Query("phone_number"),
		Email:       c.Query("email"),
	}

	// 3. ส่งต่อให้ Usecase ทำงาน
	patients, err := h.patientUsecase.SearchPatients(c.Request.Context(), hospitalID, criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. ส่งผลลัพธ์กลับเป็น JSON
	c.JSON(http.StatusOK, gin.H{
		"data": patients,
	})
}