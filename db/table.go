package db

import (
	"github.com/ongoil/agnos-test/models"
	"gorm.io/gorm"
)

func DataBaseMigration(DBConn *gorm.DB) error {
	DBConn.AutoMigrate(
		models.Hospital{},
		models.Patient{},
		models.Staff{},
	)
	return nil
}
