package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ongoil/agnos-test/handler"
	"github.com/ongoil/agnos-test/middleware"
)

func SetRouter(app *gin.Engine) {
	app.POST("/staff/create", handler.CreateStaff)
	app.POST("/staff/login", handler.Login)
	protected := app.Group("/")
	protected.Use(middleware.Auth())
	protected.POST("/patient/create", handler.CreatePatient)
	protected.GET("/patient/search", handler.SearchPatient)

	// Backward-compatible versioned routes.
	v1 := app.Group("/backend/api/v1")
	v1.POST("/register/hospital", handler.CreateHospital)
	v1.POST("/register/staff", handler.CreateStaff)
	v1.POST("/auth/login", handler.Login)
	v1.GET("/patient/search", middleware.Auth(), handler.SearchPatient)
	v1.POST("/patient/create", middleware.Auth(), handler.CreatePatient)
}
