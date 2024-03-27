package models

import (
	"time"

	"go-ranking/dao"
)

type User struct {
	Id         int    `gorm:"column:id"`
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	AddTime    int64  `gorm:"column:add_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

type UserApi struct {
	Id       int    `gorm:"column:id"`
	Username string `gorm:"column:username"`
}

func (User) TableName() string {
	return "user"
}

func GetUserById(id int) (User, error) {
	var user User
	err := dao.Db.Where("id = ?", id).First(&user).Error

	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := dao.Db.Where("username=?", username).First(&user).Error
	return user, err
}

func AddUser(username string, password string) (int, error) {
	user := User{Username: username, Password: password, AddTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	err := dao.Db.Create(&user).Error
	return user.Id, err
}

func UpdateUser(id int, username string, password string) error {
	err := dao.Db.Model(&User{}).
		Where("id = ?", id).
		Update("username", username).
		Update("password", password).Error
	return err
}

func DeleteUser(id int) error {
	err := dao.Db.Delete(&User{}, id).Error
	return err
}

func GetUserList() ([]User, error) {
	var users []User
	err := dao.Db.Find(&users).Error
	return users, err
}
