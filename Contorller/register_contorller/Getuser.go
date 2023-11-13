package registercontorller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teerapoom/Dormitory_Api/database"
	"github.com/teerapoom/Dormitory_Api/database/model"
)

// func Get_All_UserRegister(c *gin.Context) {
// 	var user_register []model.User_Register
// 	database.Db.Find(&user_register)
// 	c.JSON(http.StatusOK, gin.H{"User": user_register})
// }

func Get_All_UserRegister(c *gin.Context) {
	var user_register []model.User_Register
	err := database.Db.Preload("User_information").Find(&user_register).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "users": user_register})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	var user_register model.User_Register
	err := database.Db.Preload("User_information").First(&user_register, int(userId)).Error // Db.First จะค้นหาข้อมูลโดยอิงจาก primary key.
	// ตรวจสอบว่าข้อมูล user_register ถูกต้อง
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user_register})
}
