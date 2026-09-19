package db

import (
	"github.com/ongoil/agnos-test/models"
	"gorm.io/gorm"
)

func DataBaseMigration(dbConn *gorm.DB) error {
	return dbConn.AutoMigrate(models.Hospital{}, models.Patient{}, models.Staff{})
}
