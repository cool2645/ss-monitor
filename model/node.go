package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type Node struct {
	ID                uint `gorm:"AUTO_INCREMENT"`
	Name              string
	IPv4              string
	IPv6              string
	Ss4Json           string
	Ss6Json           string
	DomainPrefix4     string
	DomainPrefix6     string
	DomainRoot        string
	Provider          string
	DNSProvider       string
	IsCleaning        bool
	EnableWatching    bool
	EnableIPv4Testing bool
	EnableIPv6Testing bool
	EnableCleaning    bool
	OS                string
	Image             string
	DataCenter        string
	Plan              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func GetNodes(db *gorm.DB) (nodes []Node, err error) {
	err = db.Order("name asc").Find(&nodes).Error
	if err != nil {
		err = errors.Wrap(err, "GetNodes")
		return
	}
	return
}

func GetNode(db *gorm.DB, id uint) (node Node, err error) {
	err = db.Where("id = ?", id).Find(&node).Error
	if err != nil {
		err = errors.Wrap(err, "GetNode")
		return
	}
	return
}

func CreateNode(db *gorm.DB, node Node) (newNode Node, err error) {

	node.IsCleaning = false

	err = db.Create(&node).Error
	if err != nil {
		err = errors.Wrap(err, "CreateNode")
		return
	}
	newNode = node
	return
}

func UpdateNodeChangedFields(db *gorm.DB, node Node) (newNode Node, err error) {
	err = db.Model(&node).Updates(node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNode")
		return
	}
	newNode = node
	return
}

func UpdateNodeAllFields(db *gorm.DB, node Node) (newNode Node, err error) {
	err = db.Save(&node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNodeAllFields")
		return
	}
	newNode = node
	return
}

func UpdateNodeFields(db *gorm.DB, id uint, fields map[string]interface{}) (err error) {
	var node Node
	err = db.Where("id = ?", id).Find(&node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNodeFields: Find node")
		return
	}
	err = db.Model(&node).Updates(fields).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNodeFields: Update node")
		return
	}
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

func ResetNode(db *gorm.DB, id uint) (err error) {
	var node Node
	err = db.Where("id = ?", id).Find(&node).Error
	if err != nil {
		err = errors.Wrap(err, "ResetNode: Find node")
		return
	}
	err = db.Model(&node).Update("is_cleaning", false).Error
	if err != nil {
		err = errors.Wrap(err, "ResetNode: Update node")
		return
	}
	return
}

func SetNodeCleaning(db *gorm.DB, id uint) (err error) {
	var node Node
	err = db.Where("id = ?", id).Find(&node).Error
	if err != nil {
		err = errors.Wrap(err, "SetNodeCleaning: Find node")
		return
	}
	err = db.Model(&node).Update("is_cleaning", true).Error
	if err != nil {
		err = errors.Wrap(err, "SetNodeCleaning: Update node")
		return
	}
	return
}
