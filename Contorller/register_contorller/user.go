package registercontorller

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/teerapoom/Dormitory_Api/database"
	"github.com/teerapoom/Dormitory_Api/database/model"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

// json เพราะว่าจะส่งค่าเข้าไปใน api สร้างฟอม
type Register struct {
	TenantID         string                  `json:"tenantID"`
	RegisterID       string                  `json:"registerID"`
	Email            string                  `json:"email" binding:"required"`
	PassWord         string                  `json:"password" binding:"required"`
	FullName         string                  `json:"fullname" binding:"required"`
	UserRights       string                  `json:"userrights" binding:"required"` //สิทธิ์การเข้าถึง
	MobileNumber     string                  `json:"mobilenumber" binding:"required"`
	DayRegister      time.Time               `json:"dayregister"`
	User_information *model.User_information `json:"User_information"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	PassWord string `json:"password" binding:"required"`
}

func PostRegister(c *gin.Context) {
	var json Register
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExist model.User_Register
	database.Db.Where("Email = ?", json.Email).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Exist"})
		return
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.PassWord), 10) //เข้ารัหส
	json.DayRegister = time.Now()
	json.TenantID = GenerateIDUser()
	json.RegisterID = GenerateIDUser()
	for json.TenantID == json.RegisterID {
		json.RegisterID = GenerateIDUser()
	}

	json.User_information.Check_in = time.Now()
	json.User_information.Check_out = json.User_information.Check_in.Add(6 * 30 * 24 * time.Hour) // เดือน * วัน * ชั่วโมง
	userInformation := model.User_information{
		Number_Room:  json.User_information.Number_Room,
		Deposit:      json.User_information.Deposit,
		Number_phone: json.User_information.Number_phone,
		IDCard:       json.User_information.IDCard,
		Address:      json.User_information.Address,
		Check_in:     json.User_information.Check_in,
		Check_out:    json.User_information.Check_out,
	}

	User_Register := model.User_Register{
		TenantID:         json.TenantID,
		RegisterID:       json.RegisterID,
		Email:            json.Email,
		PassWord:         string(encryptedPassword),
		FullName:         json.FullName,
		UserRights:       json.UserRights,
		MobileNumber:     json.MobileNumber,
		DayRegister:      json.DayRegister,
		User_information: userInformation,
	}

	result := database.Db.Create(&User_Register)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Create Failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "User Create Success", "userId": User_Register.ID})

}

func UpdateInfo(c *gin.Context) {
	var RoomInfo model.User_Register
	id := c.Param("id")
	if err := database.Db.First(&RoomInfo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&RoomInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	var userINFO model.User_information
	database.Db.Where("number_room = ?", userINFO.Number_Room).First(&userINFO)
	if userINFO.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "NumberRoom Exist"})
		return
	}

	if err := database.Db.Save(&RoomInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":     "Ok",
		"UpdateInfo": RoomInfo})
}

// Login
func PostLogin(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExist model.User_Register
	database.Db.Where("Email = ?", json.Email).First(&userExist)
	if userExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Not have User"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.PassWord), []byte(json.PassWord)) //เปรียบเทียบรหัส
	// nil -> ที่บอกตัวแปรไม่มีค่าหรือไม่ได้ชี้ไปยังออบเจ็กต์ใดๆ
	if err == nil {
		// 200 Login สร้างสำเร็จ
		hmacSampleSecret = []byte(os.Getenv("jwt_secret_key"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"UserID": userExist.ID,
			"exp":    time.Now().Add(time.Hour * 1).Unix(), //กำหนดระยะเวลา JWT
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Login Success", "Token": tokenString})
	} else {
		// 400 login ไม่ถูกต้อง
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Login Failed"})
	}
}

// ----------------------------------
// ฟังก์ชันสร้างรหัสสุ่ม
// สร้างสตริงแบบสุ่มตัวอักษร
func GenerateIDUser() string {
	rand.Seed(time.Now().UnixNano()) //เรื่มสุ่มตัวเลขในเสี่ยววินาทีที่กด
	// สร้างตัวอักษรแรก (2 ตัว)
	firstTwoLetters := RandomString(2) // 2 ส่ง len
	// สร้างตัวเลข (4 ตัว)
	numbers := RandomNumber(4)
	return firstTwoLetters + numbers
}
func RandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) //สร้างตัวแปรสุ่ม
	var result strings.Builder                                    // typr strings.Builder ใช้เก็บข้อความ (string) เช่น a,v,d เพื่อประหยัดตัวแปร
	for i := 0; i < length; i++ {
		result.WriteByte(charset[seededRand.Intn(len(charset))]) // WriteByte เลือกจาก seededRand.Intn(len(charset) เลือกสุ่ม  charset
	}
	return result.String()
}

// สร้างตัวเลขแบบสุ่ม
func RandomNumber(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(charset[seededRand.Intn(len(charset))])
	}
	return result.String()
}

// ฟังก์ชันสร้างรหัสสุ่ม
//----------------------------------
