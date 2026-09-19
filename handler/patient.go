package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/models"
)

func SearchPatient(c *gin.Context) {
	hospitalID, err := uuid.Parse(c.GetString("hospital_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid hospital context"})
		return
	}
	var req PatientSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid search parameters"})
		return
	}
	query := db.DB.Where("hospital_id = ?", hospitalID)
	if req.NationalID != "" {
		query = query.Where("national_id = ?", req.NationalID)
	}
	if req.PassportID != "" {
		query = query.Where("passport_id = ?", req.PassportID)
	}
	if req.FirstName != "" {
		query = query.Where("(first_name_th ILIKE ? OR first_name_en ILIKE ?)", "%"+req.FirstName+"%", "%"+req.FirstName+"%")
	}
	if req.MiddleName != "" {
		query = query.Where("(middle_name_th ILIKE ? OR middle_name_en ILIKE ?)", "%"+req.MiddleName+"%", "%"+req.MiddleName+"%")
	}
	if req.LastName != "" {
		query = query.Where("(last_name_th ILIKE ? OR last_name_en ILIKE ?)", "%"+req.LastName+"%", "%"+req.LastName+"%")
	}
	if req.PhoneNumber != "" {
		query = query.Where("phone_number = ?", req.PhoneNumber)
	}
	if req.Email != "" {
		query = query.Where("email ILIKE ?", req.Email)
	}
	if req.DateOfBirth != "" {
		date, parseErr := time.Parse("2006-01-02", req.DateOfBirth)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "date_of_birth must use YYYY-MM-DD"})
			return
		}
		query = query.Where("date_of_birth = ?", date)
	}
	var patients []models.Patient
	if err := query.Order("last_name_en, last_name_th, first_name_en, first_name_th").Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to search patients"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": patients, "total": len(patients)})
}
