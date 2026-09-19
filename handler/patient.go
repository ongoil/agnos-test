package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/models"
)

func CreatePatient(c *gin.Context) {
	hospitalID, err := uuid.Parse(c.GetString("hospital_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid hospital context"})
		return
	}

	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "patient_hn is required and request body is invalid"})
		return
	}
	gender := strings.ToUpper(strings.TrimSpace(req.Gender))
	if gender != "" && gender != "M" && gender != "F" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gender must be M or F"})
		return
	}
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		parsed, parseErr := time.Parse("2006-01-02", req.DateOfBirth)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "date_of_birth must use YYYY-MM-DD"})
			return
		}
		dateOfBirth = &parsed
	}
	patient := models.Patient{
		HospitalID: hospitalID, PatientHN: strings.TrimSpace(req.PatientHN),
		NationalID: req.NationalID, PassportID: req.PassportID,
		FirstNameTH: req.FirstNameTH, MiddleNameTH: req.MiddleNameTH, LastNameTH: req.LastNameTH,
		FirstNameEN: req.FirstNameEN, MiddleNameEN: req.MiddleNameEN, LastNameEN: req.LastNameEN,
		DateOfBirth: dateOfBirth, PhoneNumber: req.PhoneNumber, Email: req.Email, Gender: gender,
	}
	if err := db.DB.Create(&patient).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"message": "patient_hn already exists in this hospital"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create patient"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": patient})
}

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
