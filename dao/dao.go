package dao

import "github.com/ydzydzydz/pmail_spam_block/model"

// ISettingDao 是设置数据访问对象的接口
type ISettingDao interface {
	// GetSetting 获取用户的设置
	GetSetting(userID int) (*model.SpamBlockSetting, error)
	// UpdateSetting 更新用户的设置
	UpdateSetting(userID int, setting *model.SpamBlockSetting) error
	// CreateSetting 创建用户的设置
	CreateSetting(setting *model.SpamBlockSetting) error
	// ExistSetting 检查用户的设置是否存在
	ExistSetting(userID int) bool
}

// IUserDao 是用户数据访问对象的接口
type IUserDao interface {
	// GetUser 获取用户
	GetUserID(account string) (int, error)
}
