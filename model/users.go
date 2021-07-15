package model

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(20);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}
