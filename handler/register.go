package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/dto"
	"github.com/ongoil/agnos-test/logger"
	"github.com/ongoil/agnos-test/models"
)

func CreateStaff(c *gin.Context) {
	var req RegisterRequest
	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(c, "400 | Bad Request : Invalid request body -> ", err, logrus.Fields{"body": req})
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "invalid request body",
			MessageTh: "ข้อมูลไม่ถูกต้อง",
		})
		return
	}

	username := strings.TrimSpace(req.Username)
	hospitalName := strings.TrimSpace(req.Hospital)

	// Validate username
	if username == "" {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "username is required",
			MessageTh: "กรุณาระบุชื่อผู้ใช้",
		})
		return
	}

	// Validate password
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "password must be at least 6 characters",
			MessageTh: "รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร",
		})
		return
	}

	// Validate hospital
	if hospitalName == "" && req.HospitalID == uuid.Nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "hospital is required",
			MessageTh: "กรุณาระบุโรงพยาบาล",
		})
		return
	}

	tx := db.DB.Begin()

	// Find hospital
	var hospital models.Hospital
	if req.HospitalID != uuid.Nil {
		if err := db.DB.Where("id = ?", req.HospitalID).First(&hospital).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, dto.Response{
					Status:    "404",
					Message:   "hospital not found",
					MessageTh: "ไม่พบโรงพยาบาล",
				})
				return
			}
			logger.LogError(c, "500 | Internal Server Error : failed to find hospital -> ", err, logrus.Fields{"hospital_id": req.HospitalID})
			c.JSON(http.StatusInternalServerError, dto.Response{
				Status:    "500",
				Message:   "failed to find hospital",
				MessageTh: "ไม่สามารถค้นหาโรงพยาบาลได้",
			})
			return
		}
	} else {
		if err := db.DB.Where("id = ?", req.HospitalID).First(&hospital).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, dto.Response{
					Status:    "404",
					Message:   "hospital not found",
					MessageTh: "ไม่พบโรงพยาบาล",
				})
				return
			}

			logger.LogError(c, "500 | Internal Server Error : failed to find hospital -> ", err, logrus.Fields{"hospital": hospitalName})
			c.JSON(http.StatusInternalServerError, dto.Response{
				Status:    "500",
				Message:   "failed to find hospital",
				MessageTh: "ไม่สามารถค้นหาโรงพยาบาลได้",
			})
			return
		}
	}

	// Check duplicate username within the same hospital
	var existing models.Staff
	if err := db.DB.Where("username = ? AND hospital_id = ?", username, hospital.Id).First(&existing).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogError(c, "500 | Internal Server Error : failed to check username -> ", err, logrus.Fields{"username": username, "hospital_id": hospital.Id})
			c.JSON(http.StatusInternalServerError, dto.Response{
				Status:    "500",
				Message:   "failed to check username",
				MessageTh: "ไม่สามารถตรวจสอบชื่อผู้ใช้ได้",
			})
			return
		}
	} else {
		c.JSON(http.StatusConflict, dto.Response{
			Status:    "409",
			Message:   "username already exists in this hospital",
			MessageTh: "ชื่อผู้ใช้นี้มีอยู่แล้วในโรงพยาบาลนี้",
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to hash password -> ", err, logrus.Fields{"username": username})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to hash password",
			MessageTh: "ไม่สามารถเข้ารหัสรหัสผ่านได้",
		})
		return
	}

	// Create staff
	staff := models.Staff{
		HospitalID:   hospital.Id,
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := tx.Create(&staff).Error; err != nil {
		tx.Rollback()
		logger.LogError(c, "500 | Internal Server Error : failed to create staff -> ", err, logrus.Fields{"username": username, "hospital_id": hospital.Id})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to create staff",
			MessageTh: "ไม่สามารถสร้างเจ้าหน้าที่ได้",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.LogError(c, "500 | Internal Server Error : ailed to commit transaction  -> ", err, logrus.Fields{"staff id": staff.Id})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "Internal Server Error - POST transactions",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
			Error:     err.Error(),
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, dto.Response{
		Status:    "201",
		Message:   "staff created successfully",
		MessageTh: "สร้างเจ้าหน้าที่สำเร็จ",
		Data: CreateStaffResponse{
			ID:         staff.Id,
			Username:   staff.Username,
			HospitalID: staff.HospitalID,
		},
	})
}
