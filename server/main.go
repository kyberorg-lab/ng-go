package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
)

func main() {
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	db, err := gorm.Open(
		"postgres",
		"host="+host+" user="+user+" password="+pass+" dbname="+dbName+" sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	app := App{
		db: db,
		r:  gin.Default(),
	}
	app.start()
}
