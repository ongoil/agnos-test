package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/dto"
	"github.com/ongoil/agnos-test/logger"
	"github.com/ongoil/agnos-test/models"
)

func CreatePatient(c *gin.Context) {
	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(c, "400 | Bad Request : Invalid request body -> ", err, logrus.Fields{"body": req})
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "patient_hn is required and request body is invalid",
			MessageTh: "กรุณาระบุ Patient HN และข้อมูลต้องถูกต้อง",
		})
		return
	}

	tx := db.DB.Begin()
	hospitalID, err := uuid.Parse(c.GetString("hospital_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Status:    "401",
			Message:   "invalid hospital context",
			MessageTh: "ข้อมูลโรงพยาบาลใน Token ไม่ถูกต้อง",
		})
		return
	}

	nationalID := strings.TrimSpace(req.NationalID)
	passportID := strings.ToUpper(strings.TrimSpace(req.PassportID))

	if nationalID == "" && passportID == "" {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "national_id or passport_id is required",
			MessageTh: "กรุณาระบุเลขบัตรประชาชนหรือเลขพาสปอร์ต",
		})
		return
	}

	query := db.DB.Where("hospital_id = ?", hospitalID)
	switch {
	case nationalID != "" && passportID != "":
		query = query.Where("(national_id = ? OR passport_id = ?)", nationalID, passportID)
	case nationalID != "":
		query = query.Where("national_id = ?", nationalID)
	default:
		query = query.Where("passport_id = ?", passportID)
	}

	var patients []models.Patient
	if err := query.Find(&patients).Error; err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to search patients -> ", err, logrus.Fields{"hospital_id": hospitalID})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to search patients",
			MessageTh: "ไม่สามารถค้นหาข้อมูลผู้ป่วยได้",
		})
		return
	}

	if len(patients) > 0 {
		c.JSON(http.StatusConflict, dto.Response{
			Status:    "409",
			Message:   "patient already exists in this hospital",
			MessageTh: "มีข้อมูลผู้ป่วยรายนี้ในโรงพยาบาลนี้แล้ว",
		})
		return
	}

	patientHN, err := generatePatientHN(db.DB, hospitalID)
	if err != nil {
		logger.LogError(c, "500 | Internal Server Error : failed to generate patient HN -> ", err, logrus.Fields{"hospital_id": hospitalID})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to generate patient_hn",
			MessageTh: "ไม่สามารถสร้าง Patient HN ได้",
		})
		return
	}

	// Validate gender
	gender := strings.ToUpper(strings.TrimSpace(req.Gender))
	if gender != "" && gender != "M" && gender != "F" {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "gender must be M or F",
			MessageTh: "Gender ต้องเป็น M หรือ F เท่านั้น",
		})
		return
	}

	// Validate date of birth
	var dateOfBirth *time.Time
	if strings.TrimSpace(req.DateOfBirth) != "" {
		parsed, parseErr := time.Parse("2006-01-02", strings.TrimSpace(req.DateOfBirth))
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, dto.Response{
				Status:    "400",
				Message:   "date_of_birth must use YYYY-MM-DD",
				MessageTh: "วันเกิดต้องอยู่ในรูปแบบ YYYY-MM-DD",
			})
			return
		}
		dateOfBirth = &parsed
	}

	// Create patient
	patient := models.Patient{
		HospitalID:   hospitalID,
		PatientHN:    patientHN,
		NationalID:   strings.TrimSpace(req.NationalID),
		PassportID:   strings.TrimSpace(req.PassportID),
		FirstNameTH:  strings.TrimSpace(req.FirstNameTH),
		MiddleNameTH: strings.TrimSpace(req.MiddleNameTH),
		LastNameTH:   strings.TrimSpace(req.LastNameTH),
		FirstNameEN:  strings.TrimSpace(req.FirstNameEN),
		MiddleNameEN: strings.TrimSpace(req.MiddleNameEN),
		LastNameEN:   strings.TrimSpace(req.LastNameEN),
		DateOfBirth:  dateOfBirth,
		PhoneNumber:  strings.TrimSpace(req.PhoneNumber),
		Email:        strings.TrimSpace(req.Email),
		Gender:       gender,
	}

	if err := tx.Create(&patient).Error; err != nil {
		tx.Rollback()
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, dto.Response{
				Status:    "409",
				Message:   "patient_hn already exists in this hospital",
				MessageTh: "Patient HN นี้มีอยู่แล้วในโรงพยาบาลนี้",
			})
			return
		}
		logger.LogError(c, "500 | Internal Server Error : failed to create patient -> ", err, logrus.Fields{"hospital_id": hospitalID, "patient_hn": patientHN})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to create patient",
			MessageTh: "ไม่สามารถสร้างข้อมูลผู้ป่วยได้",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.LogError(c, "500 | Internal Server Error : ailed to commit transaction  -> ", err, logrus.Fields{"hospital_id": hospitalID, "patient_hn": patientHN})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "Internal Server Error - POST transactions",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
			Error:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:    "200",
		Message:   "patient created successfully",
		MessageTh: "สร้างข้อมูลผู้ป่วยสำเร็จ",
		Data:      patient,
	})
}

func SearchPatient(c *gin.Context) {
	search := strings.TrimSpace(c.Param("id"))
	if search == "" {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "search is required",
			MessageTh: "กรุณาระบุคำที่ต้องการค้นหา",
		})
		return
	}

	hospitalID, err := uuid.Parse(c.GetString("hospital_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid hospital context"})
		return
	}

	query := db.DB.Where("hospital_id = ?", hospitalID)
	searchWithoutSpaces := strings.ReplaceAll(search, " ", "")
	query = query.Where(`
		STRPOS(
			REPLACE(
				CONCAT_WS(
					'|',
					patient_hn,
					national_id,
					passport_id,
					first_name_th,
					middle_name_th,
					last_name_th,
					first_name_en,
					middle_name_en,
					last_name_en
				),
				' ',
				''
			),
			?
		) > 0
	`, searchWithoutSpaces)

	var patients []models.Patient
	if err := query.Order("last_name_en, last_name_th, first_name_en, first_name_th").Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to search patients"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": patients, "total": len(patients)})
}

func generatePatientHN(dbc *gorm.DB, hospitalID uuid.UUID) (string, error) {
	var last models.Patient
	err := dbc.Unscoped(). // รวมแถวที่ soft-delete แล้ว กัน HN ซ้ำ
				Where("hospital_id = ? AND patient_hn ~ ?", hospitalID, `^HN[0-9]+$`).
				Order("LENGTH(patient_hn) DESC, patient_hn DESC").
				Limit(1).
				Take(&last).Error

	next := 1
	switch {
	case err == nil:
		n, convErr := strconv.Atoi(strings.TrimPrefix(last.PatientHN, "HN"))
		if convErr != nil {
			return "", convErr
		}
		next = n + 1
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return "", err
	}

	return fmt.Sprintf("HN%06d", next), nil
}

func GetPatient(c *gin.Context) {
	hospitalID, err := uuid.Parse(c.GetString("hospital_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid hospital context"})
		return
	}

	var patients []models.Patient
	if err := db.DB.Where("hospital_id = ?", hospitalID).Order("patient_hn").Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to search patients"})
		return
	}

	var total int64
	db.DB.Model(&models.Patient{}).Where("hospital_id = ?", hospitalID).Count(&total)

	c.JSON(http.StatusOK, dto.Response{
		Status:    "200",
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      patients,
		TotalData: total,
	})
}
