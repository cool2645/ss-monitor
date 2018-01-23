package model

import "time"

type Node struct {
	ID         uint `gorm:"AUTO_INCREMENT"`
	IPv4       string
	IPv6       string
	Ss4Json    string
	Ss6Json    string
	IsCleaning bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
