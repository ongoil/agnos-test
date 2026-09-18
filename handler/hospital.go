package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/dto"
	"github.com/ongoil/agnos-test/logger"
	"github.com/ongoil/agnos-test/models"
	"github.com/sirupsen/logrus"
)

func CreateHospital(c *gin.Context) {
	var req bodyHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(c, "400 | Bad Request : Invalid request body -> ", err, logrus.Fields{"body": req})
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:    "400",
			Message:   "invalid request body",
			MessageTh: "ข้อมูลไม่ถูกต้อง",
		})
		return
	}

	tx := db.DB.Begin()

	// Create Hospital
	newHospital := models.Hospital{
		Name: req.HospitalName,
	}

	if err := tx.Create(&newHospital).Error; err != nil {
		tx.Rollback()
		logger.LogError(c, "500 | Internal Server Error : failed to create hospital -> ", err, logrus.Fields{"hospital_name": req.HospitalName})
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:    "500",
			Message:   "failed to create hospital",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.LogError(c, "500 | Internal Server Error : ailed to commit transaction  -> ", err, logrus.Fields{"hospital_name": req.HospitalName})
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
		Message:   "Success",
		MessageTh: "สำเร็จ",
	})
}
