package models

import (
	"time"

	"gorm.io/gorm"
)

type MyClient struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:char(250);not null" json:"name"`
	Slug         string         `gorm:"type:char(100);not null" json:"client_slug"`
	IsProject    string         `gorm:"type:varchar(30);default:'0';not null" json:"is_project"`
	SelfCapture  string         `gorm:"type:char(1);default:'1';not null" json:"self_capture"`
	ClientPrefix string         `gorm:"type:char(4);not null" json:"client_prefix"`
	ClientLogo   string         `gorm:"type:char(255);default:'no-image.jpg';not null" json:"client_logo"`
	Address      string         `gorm:"type:text" json:"address"`
	PhoneNumber  string         `gorm:"type:char(50)" json:"phone_number"`
	City         string         `gorm:"type:char(50)" json:"city"`
	CreatedAt    time.Time      `gorm:"type:timestamp(0)" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp(0)" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:timestamp(0)" json:"deleted_at"`
}
