package main

import (
	"log"

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

	log.Println("server running on :5000")

	if err := app.Run(":5000"); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
