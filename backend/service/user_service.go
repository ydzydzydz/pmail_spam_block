package service

import (
	"errors"

	"github.com/Jinnrry/pmail/models"
	"xorm.io/xorm"
)

// UserService 用户服务层
type UserService struct {
	db *xorm.Engine
}

// NewUserService 创建用户服务层
func NewUserService(db *xorm.Engine) *UserService {
	return &UserService{db: db}
}

// GetUserID 获取用户ID
func (u *UserService) GetUserID(account string) (int, error) {
	var user models.User
	has, err := u.db.Where("account = ?", account).Get(&user)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errors.New("user not found")
	}
	return user.ID, nil
}
