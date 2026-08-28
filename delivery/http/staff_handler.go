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

	// กลุ่มนี้ไม่ต้องใช้ Middleware (เป็น Public API)
	public := r.Group("/staff")
	{
		public.POST("/create", handler.Create)
		public.POST("/login", handler.Login)
	}
}

// สร้าง Struct สำหรับรับ JSON Request Body
type StaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"` // รับค่า code เช่น 'HOSPITAL_A'
}

func (h *staffHandler) Create(c *gin.Context) {
	var req StaffRequest
	// c.ShouldBindJSON จะเช็คว่าส่งข้อมูลมาครบตาม tag binding:"required" หรือไม่
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

	// ถ้าล็อกอินผ่าน จะได้ JWT Token กลับมา
	token, err := h.staffUsecase.Login(c.Request.Context(), req.Username, req.Password, req.Hospital)
	if err != nil {
		// ส่ง 401 Unauthorized เสมอถ้าล็อกอินไม่ผ่าน
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}