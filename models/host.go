package models

import (
	u "apt-api/utils"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

//a struct to rep user account
type Host struct {
	gorm.Model
	Name            string `json:"name"`
	Token           string `json:"token" gorm:"unique"`
	SecurityUpdates uint   `gorm:"-" json:"securityUpdates"`
	Updates         uint   `gorm:"-" json:"updates"`
}

func (host *Host) Validate() (u.ReturnMessage, bool) {
	return u.Message(false, "Validation passed"), true
}

func (host *Host) Create() u.ReturnMessage {
	host.Token = GenerateHostToken()
	GetDB().Create(host)

	response := u.Message(true, "Host has been created")
	response["host"] = host
	return response
}

func GenerateHostToken() string {
	return uuid.New().String()
}

func (host *Host) Delete() u.ReturnMessage {
	GetDB().Delete(&host)

	return u.Message(true, "host deleted")
}

func GetHosts() []*Host {
	hosts := make([]*Host, 0)
	err := GetDB().Table("hosts").Find(&hosts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return hosts
}

func GetHostMatchingToken(token string) (*Host, error) {
	var host Host
	err := GetDB().Where("token = ?", token).First(&host).Error
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (h *Host) CountUpdates() uint {
	var count uint
	GetDB().Model(&Update{}).Where("host_id = ?", h.ID).Count(&count)
	return count
}

func (h *Host) CountSecurityUpdates() uint {
	var count uint
	GetDB().Model(&Update{}).Where("host_id = ? AND security = true", h.ID).Count(&count)
	return count
}

func GetHost(id uint) *Host {
	host := &Host{}
	db.First(&host, id)

	return host
}
