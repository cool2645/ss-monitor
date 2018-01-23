package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

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

func NodesStatus() {

}

func CreateNode(db *gorm.DB, node Node) (newNode Node, err error) {
	// Default Value
	node.IsCleaning = false

	err = db.Create(&node).Error
	if err != nil {
		err = errors.Wrap(err, "CreateNode")
		return
	}
	newNode = node
	return
}

func UpdateNode(db *gorm.DB, node Node) (newNode Node, err error) {
	err = db.Model(&node).Updates(node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNode")
		return
	}
	newNode = node
	return
}

func DeleteNode(db *gorm.DB, id uint) (err error) {
	err = db.Delete(Node{}, "id = ?", id).Error
	if err != nil {
		err = errors.Wrap(err, "DeleteNode")
		return
	}
	return
}
