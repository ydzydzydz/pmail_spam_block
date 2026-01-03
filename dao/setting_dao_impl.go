package dao

import (
	"github.com/ydzydzydz/pmail_spam_block/db/repository"
	"github.com/ydzydzydz/pmail_spam_block/model"
	"xorm.io/xorm"
)

// SettingDaoImpl 实现了 ISettingDao 接口
type SettingDaoImpl struct {
	db *xorm.Engine
}

var _ ISettingDao = (*SettingDaoImpl)(nil)

// NewSettingDaoImpl 创建一个新的 SettingDaoImpl 实例
func NewSettingDaoImpl(db *xorm.Engine) *SettingDaoImpl {
	return &SettingDaoImpl{db: db}
}

// GetSetting 获取用户的设置
func (d *SettingDaoImpl) GetSetting(userID int) (*model.SpamBlockSetting, error) {
	settingRepo := repository.NewRepository[model.SpamBlockSetting](d.db)
	return settingRepo.FindOne(userID)
}

// UpdateSetting 更新用户的设置
func (d *SettingDaoImpl) UpdateSetting(userID int, setting *model.SpamBlockSetting) error {
	settingRepo := repository.NewRepository[model.SpamBlockSetting](d.db)
	return settingRepo.Update(userID, setting)
}

// CreateSetting 创建用户的设置
func (d *SettingDaoImpl) CreateSetting(setting *model.SpamBlockSetting) error {
	settingRepo := repository.NewRepository[model.SpamBlockSetting](d.db)
	return settingRepo.Create(setting)
}

// ExistSetting 检查用户的设置是否存在
func (d *SettingDaoImpl) ExistSetting(userID int) bool {
	settingRepo := repository.NewRepository[model.SpamBlockSetting](d.db)
	return settingRepo.Exist(userID)
}
