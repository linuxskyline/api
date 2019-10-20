package models

import (
	u "apt-api/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

//a struct to rep user account
type Update struct {
	gorm.Model
	HostId         uint   `json:"host_id" gorm:"unique_index:host_package_name"`
	PackageName    string `json:"packageName" gorm:"unique_index:host_package_name"`
	CurrentVersion string `json:"currentVersion"`
	NewVersion     string `json:"newVersion"`
	Security       bool   `json:"security"`
}

func (update *Update) Create() u.ReturnMessage {
	currentUpdate := Update{}

	if err := GetDB().Where("package_name = ? AND host_id = ?", update.PackageName, update.HostId).First(&currentUpdate).Error; err != nil {
		GetDB().Create(update)

		response := u.Message(true, fmt.Sprintf("Update has been created %s", err))
		response["update"] = update
		return response
	}

	if update.NewVersion == currentUpdate.NewVersion {
		response := u.Message(true, "Update already up to date")
		response["update"] = currentUpdate
		return response
	}

	currentUpdate.NewVersion = update.NewVersion
	currentUpdate.Security = update.Security
	currentUpdate.CurrentVersion = update.CurrentVersion

	GetDB().Save(&currentUpdate)

	response := u.Message(true, "Updated existing update")
	response["update"] = currentUpdate
	return response
}

func GetUpdates(host uint) []*Update {

	updates := make([]*Update, 0)
	err := GetDB().Table("updates").Where("host_id = ?", host).Find(&updates).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return updates
}
