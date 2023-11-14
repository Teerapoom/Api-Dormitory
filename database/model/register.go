package model

import (
	"time"

	"gorm.io/gorm"
)

// ORM
type User_Register struct {
	gorm.Model
	Email            string
	PassWord         string
	FullName         string
	UserRights       string
	MobileNumber     string
	DayRegister      time.Time
	TenantID         string
	RegisterID       string
	User_information User_information `gorm:"foreignKey:UserID"`
}

type User_information struct {
	gorm.Model
	UserID       uint
	Number_Room  string
	Deposit      float64
	Number_phone string
	IDCard       string
	Address      string
	Check_in     time.Time
	Check_out    time.Time
}

/*ห้องพัก{
	เลขห้อง
	จำนวนคนเข้าพัก
	หน่วยไฟ

}
*/
