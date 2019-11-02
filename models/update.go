package models

import (
	u "apt-api/utils"
	"fmt"
	"time"

)

//a struct to rep user account
type Update struct {
	ID                  int       `gorm:"primary_key" sql:"AUTO_INCREMENT" json:"id" jsonapi:"primary,updates"`
	CreatedAt           time.Time `sql:"DEFAULT:current_timestamp" json:"created_at" jsonapi:"attr,created_at"`
	UpdatedAt           time.Time `json:"updated_at" jsonapi:"attr,updated_at"`
	HostId         uint   `json:"host_id" gorm:"unique_index:host_package_name" jsonapi:"attr,hostId"`
	PackageName    string `json:"packageName" gorm:"unique_index:host_package_name" jsonapi:"attr,packageName"`
	CurrentVersion string `json:"currentVersion" jsonapi:"attr,currentVersion"`
	NewVersion     string `json:"newVersion" jsonapi:"attr,newVersion"`
	Security       bool   `json:"security" jsonapi:"attr,security"`
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

func (update *Update) Delete() u.ReturnMessage {
	GetDB().Delete(&update)

	return u.Message(true, "Update deleted")
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

func GetUpdate(id uint) *Update {
	update := &Update{}
	db.First(&update, id)

	return update
}
