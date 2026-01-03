package service

import (
	"github.com/ydzydzydz/pmail_spam_block/dao"
	"github.com/ydzydzydz/pmail_spam_block/model"
)

// SettingService 设置服务
type SettingService struct {
	dao dao.ISettingDao
}

const (
	DefaultApiUrl    = ""    // 默认api url, 为空时, 则不发送
	DefaultTimeout   = 50000 // 默认超时时间
	DefaultThreshold = 0.2   // 默认阈值
)

// NewSettingService 创建设置服务实例
func NewSettingService(dao dao.ISettingDao) *SettingService {
	return &SettingService{dao: dao}
}

// GetSetting 获取设置
// 如果不存在, 则创建默认设置
func (s *SettingService) GetSetting(userID int) (*model.SpamBlockSetting, error) {
	has := s.dao.ExistSetting(userID)
	if !has {
		if err := s.CreateDefaultSetting(userID); err != nil {
			return nil, err
		}
	}
	return s.dao.GetSetting(userID)
}

// UpdateSetting 更新设置
// 如果不存在, 则创建默认设置
func (s *SettingService) UpdateSetting(userID int, setting *model.SpamBlockSetting) error {
	has := s.dao.ExistSetting(setting.UserID)
	if !has {
		if err := s.CreateDefaultSetting(setting.UserID); err != nil {
			return err
		}
	}
	return s.dao.UpdateSetting(userID, setting)
}

// CreateDefaultSetting 创建默认设置
func (s *SettingService) CreateDefaultSetting(userID int) error {
	setting := &model.SpamBlockSetting{
		UserID:    userID,
		ApiUrl:    DefaultApiUrl,
		Timeout:   DefaultTimeout,
		Threshold: DefaultThreshold,
	}
	return s.dao.CreateSetting(setting)
}
