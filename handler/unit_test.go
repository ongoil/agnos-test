package handler

import (
	"net/http"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetStaff() ([]models.Staff, error) {
	var staff []models.Staff

	err := db.DB.Find(&staff).Error

	return staff, err
}

func TestLogin_WrongPassword(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	db.DB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	password, _ := bcrypt.GenerateFromPassword(
		[]byte("123456"),
		bcrypt.DefaultCost,
	)

	mock.ExpectQuery(`SELECT .* FROM "staffs"`).
		WithArgs("test", 1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"hospital_id",
				"username",
				"password_hash",
			}).AddRow(
				"11111111-1111-1111-1111-111111111111",
				"22222222-2222-2222-2222-222222222222",
				"test",
				string(password),
			),
		)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/login", Login)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{"username":"test","password":"wrong"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
