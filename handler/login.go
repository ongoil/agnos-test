package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/dto"
	"github.com/ongoil/agnos-test/logger"
	"github.com/ongoil/agnos-test/middleware"
	"github.com/ongoil/agnos-test/models"
)

func Login(c *gin.Context) {
	var req LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(c, "400 | Bad Request : Invalid request body -> ", err, logrus.Fields{"body": req})
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "username and password are required",
			MessageTh: "กรุณาระบุชื่อผู้ใช้และรหัสผ่าน",
		})
		return
	}

	username := strings.TrimSpace(req.Username)

	var staff models.Staff
	query := db.DB.Where("username = ?", username)
	// Find staff
	if err := query.First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Status:    "401",
				Message:   "invalid username, password",
				MessageTh: "ชื่อผู้ใช้ หรือรหัสผ่านไม่ถูกต้อง",
			})
			return
		}
		logger.LogError(c, "500 | Internal Server Error : failed to find staff -> ", err, logrus.Fields{"username": username})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to login",
			MessageTh: "ไม่สามารถเข้าสู่ระบบได้",
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Status:    "401",
			Message:   "invalid username, password",
			MessageTh: "ชื่อผู้ใช้ หรือรหัสผ่านไม่ถูกต้อง",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(staff.Id, staff.HospitalID)
	if err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to generate token -> ", err, logrus.Fields{"staff id": staff.Id})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to generate token",
			MessageTh: "ไม่สามารถสร้าง Token ได้",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, dto.Response{
		Status:    "200",
		Message:   "login successfully",
		MessageTh: "เข้าสู่ระบบสำเร็จ",
		Data: LoginStaffResponse{
			Token: token,
		},
	})
}
