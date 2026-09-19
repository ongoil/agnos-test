package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ongoil/agnos-test/handler"
	"github.com/ongoil/agnos-test/middleware"
)

func SetRouter(app *gin.Engine) {
	// Public routes
	app.POST("/staff/create", handler.CreateStaff) // create ข้อมูล staff

	v1 := app.Group("/backend/api/v1")
	// Public v1 routes
	v1.POST("/register/hospital", handler.CreateHospital) // สรา้งข้อมูลโรงบาล
	v1.POST("/register/staff", handler.CreateStaff)       // สรา้งข้อมูล staff
	v1.POST("/auth/login", handler.Login)                 // login

	// Protected v1 routes
	v1Protected := v1.Group("/")
	v1Protected.Use(middleware.Auth())
	v1Protected.POST("/patient/create", handler.CreatePatient) // สรา้งข้อมูลข้อมูลผู้ป่วย patient
	v1Protected.GET("/patient/search", handler.SearchPatient)  // ค้นหาข้อมูลผู้ป่วย patient

}
