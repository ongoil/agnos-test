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
	"github.com/ongoil/agnos-test/models"
)

func CreateStaff(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username, password (minimum 6 characters), and hospital are required"})
		return
	}
	if strings.TrimSpace(req.Hospital) == "" && req.HospitalID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "hospital is required"})
		return
	}
	var hospital models.Hospital
	query := db.DB
	if req.HospitalID != uuid.Nil {
		if err := query.First(&hospital, "id = ?", req.HospitalID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "hospital not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to find hospital"})
			return
		}
	} else if err := query.Where("name = ?", strings.TrimSpace(req.Hospital)).First(&hospital).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "hospital not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to find hospital"})
		return
	}
	var existing models.Staff
	if err := query.Where("username = ?", strings.TrimSpace(req.Username)).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "username already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to check username"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to hash password"})
		return
	}
	staff := models.Staff{HospitalID: hospital.Id, Username: strings.TrimSpace(req.Username), PasswordHash: string(hash)}
	if err := query.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create staff"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": CreateStaffResponse{ID: staff.Id, Username: staff.Username, HospitalID: staff.HospitalID}})
}
