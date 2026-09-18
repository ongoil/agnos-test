package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ongoil/agnos-test/handler"
	"github.com/ongoil/agnos-test/middleware"
)

func TestAdd(t *testing.T) {
	result := Add(2, 3)

	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func Add(a int, b int) int {
	return a + b
}

func SetRouter(app *gin.Engine) {
	v1 := app.Group("/backend/api/v1")

	register := v1.Group("/register")
	register.POST("/hospital", handler.CreateHospital)
	register.POST("/staff", handler.CreateStaff)

	auth := v1.Group("/auth")
	auth.POST("/login", handler.Login)

	protected := v1.Group("")
	protected.Use(middleware.Auth())

	staff := protected.Group("/staff")
	staff.GET("/show", handler.GetStaffData)

	// patient := v1.Group("/patients")
	// patient.GET("/patients/search", handler.SearchPatient)
}
