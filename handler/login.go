package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/middleware"
	"github.com/ongoil/agnos-test/models"
)

func Login(c *gin.Context) {
	var req LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username and password are required"})
		return
	}
	var staff models.Staff
	query := db.DB.Where("username = ?", strings.TrimSpace(req.Username))
	if req.HospitalID != uuid.Nil {
		query = query.Where("hospital_id = ?", req.HospitalID)
	} else if strings.TrimSpace(req.Hospital) != "" {
		query = query.Joins("JOIN hospitals ON hospitals.id = staffs.hospital_id").Where("hospitals.name = ?", strings.TrimSpace(req.Hospital))
	}
	if err := query.First(&staff).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username, password, or hospital"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to login"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username, password, or hospital"})
		return
	}
	token, err := middleware.GenerateToken(staff.Id, staff.HospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": LoginStaffResponse{Token: token}})
}
