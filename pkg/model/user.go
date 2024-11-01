package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	UserName string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
	Phone    string `gorm:"not null;unique"`
}
