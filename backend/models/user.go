package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//UserName string `json:"username"`
	UserEmail string `json:"useremail"`
	UserId    string `json:"userId"`
	Password  string `json:"password"`
}
