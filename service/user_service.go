package service

import "github.com/ydzydzydz/pmail_spam_block/dao"

// UserService 用户服务层
type UserService struct {
	userDao dao.IUserDao
}

// NewUserService 创建用户服务层
func NewUserService(userDao dao.IUserDao) *UserService {
	return &UserService{userDao: userDao}
}

// GetUserID 获取用户ID
func (u *UserService) GetUserID(account string) (int, error) {
	return u.userDao.GetUserID(account)
}
