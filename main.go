package main

import (
	"database/sql"
	"fmt"
	"ipr/roles"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	db := openDb()
	roleRepo := roles.NewRoleRepo(db)
	roleValidator := roles.NewRoleValidator()
	roleService := roles.NewRoleService(roleRepo, roleValidator)
	roleController := roles.NewRoleController(roleService)

	r := gin.New()
	r.Use(ErrorHandler)
	r.GET("/api/roles", roleController.GetAllRoles)
	r.POST("/api/roles/create", roleController.CreateRole)
	r.POST("/api/roles/update", roleController.UpdateRole)
	r.POST("/api/roles/delete", roleController.DeleteRole)
	r.Run(":7777")
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	if len(c.Errors) == 1 {
		err := c.Errors[0]
		if validationErr, ok := err.Err.(*roles.ValidationErrors); ok {
			c.JSON(http.StatusOK, validationErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"title": err.Error()})
		return
	}

	// тк мб несколько ошибок
	errs := []string{}
	for _, err := range c.Errors {
		errs = append(errs, err.Error())
	}
	errMsg := strings.Join(errs, "; ")
	c.JSON(http.StatusInternalServerError, gin.H{"title": errMsg})
}

func openDb() *sql.DB {
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	connStr := fmt.Sprintf("user=%s password=%s dbname=tracker_go sslmode=disable", user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")
	return db
}
