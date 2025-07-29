package repository

import (
	"AUV/config"
	"AUV/db"
	"AUV/models"
)

type UserRepository struct{}

var UserRepo = &UserRepository{}

func (r *UserRepository) CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}

func (r *UserRepository) GetUserByID(userId string) (*models.User, error) {
	var user models.User
	if err := db.DB.First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUserActive() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Where("username != ? and is_active = true", config.Cfg.Admin.Username).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetAllUserInActive() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Where("username != ? and is_active = false", config.Cfg.Admin.Username).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUserAllInfo(user *models.User) error {
	return db.DB.Save(user).Error
}

func (r *UserRepository) UpdateUserOtherInfo(userId string, user *models.User) error {
	return db.DB.Model(&models.User{}).Omit("id", "password", "created_at", "last_login", "is_active").Where("id = ?", userId).Updates(user).Error
}

func (r *UserRepository) UpdateUserStatus(userId string, isActive bool) error {
	return db.DB.Model(&models.User{}).Where("id = ?", userId).Update("is_active", isActive).Error
}

func (r *UserRepository) DeleteUserByID(userId string) error {
	return db.DB.Unscoped().Where("id = ?", userId).Delete(&models.User{}).Error
}
