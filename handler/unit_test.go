package handler

import (
	"net/http"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/ongoil/agnos-test/db"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLogin_WrongPassword(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	t.Cleanup(func() {
		sqlDB.Close()
		db.DB = nil
	})

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

func TestCreateStaff_RequiresHospital(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/staff/create", CreateStaff)
	req := httptest.NewRequest(http.MethodPost, "/staff/create",
		strings.NewReader(`{"username":"alice","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSearchPatient_ByID(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	t.Cleanup(func() {
		sqlDB.Close()
		db.DB = nil
	})

	db.DB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	mock.ExpectQuery(`SELECT .* FROM "patients"`).
		WithArgs(
			"22222222-2222-2222-2222-222222222222",
			"11111111-1111-1111-1111-111111111111",
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"hospital_id",
			"patient_hn",
		}).AddRow(
			"11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222",
			"HN001",
		))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/patient/search/:id", func(c *gin.Context) {
		c.Set("hospital_id", "22222222-2222-2222-2222-222222222222")
		SearchPatient(c)
	})
	req := httptest.NewRequest(
		http.MethodGet,
		"/patient/search/11111111-1111-1111-1111-111111111111",
		nil,
	)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePatient_InvalidGender(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	t.Cleanup(func() {
		sqlDB.Close()
		db.DB = nil
	})

	db.DB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	mock.ExpectBegin()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/patient/create", func(c *gin.Context) {
		c.Set("hospital_id", "22222222-2222-2222-2222-222222222222")
		CreatePatient(c)
	})
	req := httptest.NewRequest(http.MethodPost, "/patient/create",
		strings.NewReader(`{"patient_hn":"HN001","gender":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
