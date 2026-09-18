package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/dto"
	"github.com/ongoil/agnos-test/logger"
	"github.com/ongoil/agnos-test/models"
)

func CreateStaff(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(c, "400 | Bad Request : Invalid request body -> ", err, logrus.Fields{"body": req})
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "invalid request body",
			MessageTh: "ข้อมูลไม่ถูกต้อง",
		})
		return
	}

	// Find hospital
	var hospital models.Hospital
	if err := db.DB.Where("id = ?", req.HospitalID).First(&hospital).Error; err != nil {
		logger.LogError(c, "404 | Not Found : hospital not found -> ", err, logrus.Fields{"hospital_id": req.HospitalID})
		c.JSON(http.StatusNotFound, dto.Response{
			Status:    "404",
			Message:   "hospital not found",
			MessageTh: "ไม่พบข้อมูล",
		})
		return
	}

	// Check duplicate username
	var existingStaff models.Staff
	if err := db.DB.Where("username = ?", req.Username).First(&existingStaff).Error; err == nil {
		logger.LogError(c, "409 | Conflict : username already exists -> ", err, logrus.Fields{"username": req.Username})
		c.JSON(http.StatusConflict, dto.Response{
			Status:    "409",
			Message:   "username already exists",
			MessageTh: "ชื่อผู้ใช้งานนี้ถูกใช้งานแล้ว",
		})
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to hash password -> ", err, logrus.Fields{"username": req.Username})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to hash password",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
		})
		return
	}

	// Create staff
	staff := models.Staff{
		HospitalID:   hospital.Id,
		Username:     req.Username,
		PasswordHash: string(passwordHash),
	}

	if err := db.DB.Create(&staff).Error; err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to create staff -> ", err, logrus.Fields{"username": req.Username})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to create staff",
			MessageTh: "บันทึกข้อมูลไม่สำเร็จ",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:    "200",
		Message:   "Success",
		MessageTh: "สำเร็จ",
	})
}

func GetStaffData(c *gin.Context) {
	var staff []models.Staff
	if err := db.DB.Preload("Hospital").Find(&staff).Error; err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to get staff -> ", err, logrus.Fields{})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to get staff",
			MessageTh: "เกิดข้อผิดพลาดในการดึงข้อมูลเจ้าหน้าที่",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:    "200",
		Message:   "Success",
		MessageTh: "สำเร็จ",
		Data:      staff,
	})
}
