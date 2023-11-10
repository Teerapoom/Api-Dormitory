package model

import (
	"time"

	"gorm.io/gorm"
)

// ORM
type User_Register struct {
	gorm.Model
	IDUser       string
	Email        string
	PassWord     string
	FullName     string
	UserRights   string
	MobileNumber string
	DayRegister  time.Time
}

// Table user_information {
// 	id int [PK]
// 	created_at timestamp
// 	updated_at timestamp
// 	deleted_at timestamp
// 	id_information int
// 	Numder_Room string
// 	Deposit string //เงินมัดจำ
// 	Numder_phone string
// 	IDCard string // บัตรประชาชน
// 	Copy_IDCard string // สำเนาบัตรประชาชน
// 	Address string
// 	Check_in string
// 	Check_out string
//   }
