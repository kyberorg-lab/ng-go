package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"path"
	"path/filepath"
)

type student struct {
	ID   string `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type App struct {
	db *gorm.DB
	r  *gin.Engine
}

func (a *App) start() {
	a.db.AutoMigrate(&student{})
	a.r.GET("/students", a.getAllStudents)
	a.r.POST("/students", a.addStudent)
	a.r.PUT("/students/:id", a.updateStudent)
	a.r.DELETE("/students/:id", a.deleteStudent)

	a.r.NoRoute(a.serveStatic)
	log.Fatal(a.r.Run(":8080"))
}

func (a *App) serveStatic(c *gin.Context) {
	dir, file := path.Split(c.Request.RequestURI)
	ext := filepath.Ext(file)
	if file == "" || ext == "" {
		c.File("./webapp/dist/webapp/index.html")
	} else {
		c.File("./webapp/dist/webapp/" + path.Join(dir, file))
	}
}

func (a *App) getAllStudents(context *gin.Context) {
	context.Header("Content-Type", "application/json")
	var all []student
	err := a.db.Find(&all).Error
	if err != nil {

		sendErr(context, http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, all)
}

func (a *App) addStudent(context *gin.Context) {
	context.Header("Content-Type", "application/json")
	var s student
	err := context.ShouldBindJSON(&s)
	if err != nil {
		sendErr(context, http.StatusBadRequest, err.Error())
		return
	}
	s.ID = uuid.New().String()
	err = a.db.Save(&s).Error
	if err != nil {
		sendErr(context, http.StatusInternalServerError, err.Error())
	} else {
		context.JSON(http.StatusCreated, "")
	}
}

func (a *App) updateStudent(context *gin.Context) {
	context.Header("Content-Type", "application/json")
	var s student
	err := context.ShouldBindJSON(&s)
	if err != nil {
		sendErr(context, http.StatusBadRequest, err.Error())
		return
	}
	s.ID = context.Param("id")
	err = a.db.Save(&s).Error
	if err != nil {
		sendErr(context, http.StatusInternalServerError, err.Error())
	}
}

func (a *App) deleteStudent(context *gin.Context) {
	context.Header("Content-Type", "application/json")
	id := context.Param("id")
	err := a.db.Unscoped().Delete(student{ID: id}).Error
	if err != nil {
		sendErr(context, http.StatusInternalServerError, err.Error())
	}
}

func sendErr(context *gin.Context, code int, message string) {
	context.JSON(code, gin.H{"error": message})
}
