package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db}
}

func (u *UserRepository) AddUser(user model.User) error {
	if err := u.db.Create(&user).Error; err != nil {
		// return any error will rollback
		return err
	}
	return nil // TODO: replace this
}

func (u *UserRepository) UserAvail(cred model.User) error {
	var data model.User
	result := u.db.Table("users").Select("*").Where("username = ?", cred.Username).Where("password = ?", cred.Password).Scan(&data)
	if result.Error != nil {
		return result.Error
	}

	if data == (model.User{}) {
		return fmt.Errorf("user tidak available")
	}
	return nil // TODO: replace this
}

func (u *UserRepository) CheckPassLength(pass string) bool {
	if len(pass) <= 5 {
		return true
	}

	return false
}

func (u *UserRepository) CheckPassAlphabet(pass string) bool {
	for _, charVariable := range pass {
		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') {
			return false
		}
	}
	return true
}
