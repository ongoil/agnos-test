package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ongoil/agnos-test/db"
	"github.com/ongoil/agnos-test/router"
)

func Add(a int, b int) int {
	return a + b
}

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.DataBaseMigration(dbConn); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	app := gin.Default()
	router.SetRouter(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("server running on :" + port)

	if err := app.Run(":" + port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
