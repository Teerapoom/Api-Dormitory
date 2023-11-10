package registercontorller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teerapoom/Dormitory_Api/database"
	"github.com/teerapoom/Dormitory_Api/database/model"
)

func Get_All_UserRegister(c *gin.Context) {
	var user_register []model.User_Register
	database.Db.Find(&user_register)
	c.JSON(http.StatusOK, gin.H{"User": user_register})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	var user_register model.User_Register
	database.Db.First(&user_register, userId) // Db.First จะค้นหาข้อมูลโดยอิงจาก primary key.
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user_register})
}
