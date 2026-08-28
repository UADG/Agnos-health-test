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

	protected := r.Group("/patient")
	protected.Use(jwtMiddleware)
	{
		protected.GET("/search", handler.Search)
	}
}

func (h *patientHandler) Search(c *gin.Context) {
	hospitalID := c.MustGet("hospital_id").(string)

	criteria := domain.PatientSearchCriteria{
		NationalID:  c.Query("national_id"),
		PassportID:  c.Query("passport_id"),
		FirstName:   c.Query("first_name"),
		LastName:    c.Query("last_name"),
		DateOfBirth: c.Query("date_of_birth"),
		PhoneNumber: c.Query("phone_number"),
		Email:       c.Query("email"),
	}

	patients, err := h.patientUsecase.SearchPatients(c.Request.Context(), hospitalID, criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": patients,
	})
}