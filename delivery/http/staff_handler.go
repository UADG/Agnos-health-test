package http

import (
	"net/http"

	"agnos-assignment/domain"
	"github.com/gin-gonic/gin"
)

type staffHandler struct {
	staffUsecase domain.StaffUsecase
}

func NewStaffHandler(r *gin.Engine, su domain.StaffUsecase) {
	handler := &staffHandler{staffUsecase: su}

	public := r.Group("/staff")
	{
		public.POST("/create", handler.Create)
		public.POST("/login", handler.Login)
	}
}

type StaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"` 
}

func (h *staffHandler) Create(c *gin.Context) {
	var req StaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.staffUsecase.CreateStaff(c.Request.Context(), req.Username, req.Password, req.Hospital)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "staff created successfully"})
}

func (h *staffHandler) Login(c *gin.Context) {
	var req StaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.staffUsecase.Login(c.Request.Context(), req.Username, req.Password, req.Hospital)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}