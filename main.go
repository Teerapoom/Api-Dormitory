package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	registercontorller "github.com/teerapoom/Dormitory_Api/Contorller/register_contorller"
	"github.com/teerapoom/Dormitory_Api/database"
	"github.com/teerapoom/Dormitory_Api/middleware"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
	database.IntnDb()
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/registerUser", registercontorller.PostRegister)
	r.PUT("/updateInfo/:id", registercontorller.UpdateInfo)
	r.POST("/login", registercontorller.PostLogin)
	r.GET("/GetAllTest", registercontorller.Get_All_UserRegister)
	authorized := r.Group("/users", middleware.JWTAuthen()) // มีการเรียนmiddlewareทุกครั้ง
	authorized.GET("/GetAll_UserRegister", registercontorller.Get_All_UserRegister)
	authorized.GET("/Get_profile", registercontorller.Profile)
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

/*
	เวลา login เเล้ว เป็น Admin จะให้ token
	แต่ถ้าไม่เป็น Admin นั้นคือ user จะไม่ให้
*/
