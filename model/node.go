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

func UpdateNode(db *gorm.DB, node Node) (newNode Node, err error) {
	err = db.Model(&node).Updates(node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNode")
		return
	}
	newNode = node
	return
}

func UpdateNodeSetEnable(db *gorm.DB, id uint, task string) (err error) {
	var node Node
	err = db.Where("id = ?", id).Find(&node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNodeSetEnable: Find node")
		return
	}
	switch task {
	case "watching":
		node.EnableWatching = true
	case "ipv4testing":
		node.EnableIPv4Testing = true
	case "ipv6testing":
		node.EnableIPv6Testing = true
	case "cleaning":
		node.EnableCleaning = true
	default:
		err = errors.New("Unknown task type")
		err = errors.Wrap(err, "UpdateNodeSetEnable")
		return
	}
	err = db.Model(&node).Updates(node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNodeSetEnable: Update node")
		return
	}
	return
}

func UpdateNodeSetDisable(db *gorm.DB, id uint, task string) (err error) {
	var node Node
	err = db.Where("id = ?", id).Find(&node).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateNodeSetDisable: Find node")
		return
	}
	switch task {
	case "watching":
		err = db.Model(&node).Update("enable_watching", false).Error
	case "ipv4testing":
		err = db.Model(&node).Update("enable_ipv4_testing", false).Error
	case "ipv6testing":
		err = db.Model(&node).Update("enable_ipv6_testing", false).Error
	case "cleaning":
		err = db.Model(&node).Update("cleaning", false).Error
	default:
		err = errors.New("Unknown task type")
		err = errors.Wrap(err, "UpdateNodeSetDisable")
		return
	}
	if err != nil {
		err = errors.Wrap(err, "UpdateNodeSetDisable: Update node")
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
