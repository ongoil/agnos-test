package handler

import (
	"errors"
	"net/http"

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
			Message:   "invalid request body",
			MessageTh: "ข้อมูลไม่ถูกต้อง",
		})
		return
	}

	// Find staff
	var staff models.Staff
	err := db.DB.Where("username = ?", req.Username).First(&staff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Status:    "401",
			Message:   "invalid username or password",
			MessageTh: "ชื่อผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง",
		})
		return
	}

	if err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to find staff -> ", err, logrus.Fields{"username": req.Username})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to login",
			MessageTh: "เกิดข้อผิดพลาดในการเข้าสู่ระบบ",
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Status:    "401",
			Message:   "invalid username or password",
			MessageTh: "ชื่อผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง",
		})
		return
	}

	// Generate JWT
	token, err := middleware.GenerateToken(
		staff.Id,
		staff.HospitalID,
	)

	if err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to generate token -> ", err, logrus.Fields{"username": req.Username})

		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to generate token",
			MessageTh: "เกิดข้อผิดพลาดในการสร้าง Token",
		})
		return
	}

	// TODO: Generate JWT
	c.JSON(http.StatusOK, dto.Response{
		Status:    "200",
		Message:   "login successful",
		MessageTh: "เข้าสู่ระบบสำเร็จ",
		Data: LoginStaffResponse{
			Token: token,
		},
	})
}
