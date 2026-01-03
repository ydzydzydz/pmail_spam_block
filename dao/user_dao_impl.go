package dao

import (
	"errors"

	"github.com/Jinnrry/pmail/models"
	"xorm.io/xorm"
)

// UserDaoImpl 用户数据访问层实现
type UserDaoImpl struct {
	db *xorm.Engine
}

var _ IUserDao = (*UserDaoImpl)(nil)

// GetUserID 获取用户ID
func (u *UserDaoImpl) GetUserID(account string) (int, error) {
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

// NewUserDaoImpl 创建用户数据访问层实现
func NewUserDaoImpl(db *xorm.Engine) *UserDaoImpl {
	return &UserDaoImpl{db: db}
}
